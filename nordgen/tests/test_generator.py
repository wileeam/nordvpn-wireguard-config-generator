"""Tests for the Generator class"""
import pytest
from math import radians, sin, cos, asin, sqrt
from nord_config_generator.generator import Generator
from nord_config_generator.client import NordClient
from nord_config_generator.ui import ConsoleManager
from nord_config_generator.models import Server


class TestGenerator:
    """Test suite for Generator class"""

    @pytest.fixture
    def generator(self):
        """Create a generator instance for testing"""
        client = NordClient()
        ui = ConsoleManager()
        return Generator(client, ui)

    def test_check_version_valid(self, generator):
        """Test version checking with valid versions"""
        assert generator._check_version("2.1.0") == True
        assert generator._check_version("2.2.0") == True
        assert generator._check_version("3.0.0") == True
        assert generator._check_version("3.1.5") == True
        assert generator._check_version("10.5.2") == True

    def test_check_version_invalid(self, generator):
        """Test version checking with invalid versions"""
        assert generator._check_version("2.0.0") == False
        assert generator._check_version("1.9.9") == False
        assert generator._check_version("0.5.0") == False
        assert generator._check_version("1.0") == False
        assert generator._check_version("invalid") == False
        assert generator._check_version("") == False
        assert generator._check_version("2") == False
        assert generator._check_version("abc.def.ghi") == False

    def test_check_version_edge_cases(self, generator):
        """Test version checking edge cases"""
        assert generator._check_version("2.1") == True
        assert generator._check_version("2.0") == False
        assert generator._check_version("3") == False  # Less than 3 chars

    def test_haversine_same_location(self, generator):
        """Test haversine distance for same location"""
        distance = generator._haversine(40.7128, -74.0060, 40.7128, -74.0060)
        assert distance == 0.0

    def test_haversine_known_distance(self, generator):
        """Test haversine distance with known locations"""
        # New York to Los Angeles (approximate distance ~3944 km)
        ny_lat, ny_lon = 40.7128, -74.0060
        la_lat, la_lon = 34.0522, -118.2437

        distance = generator._haversine(ny_lat, ny_lon, la_lat, la_lon)

        # Allow 1% tolerance
        assert 3900 < distance < 4000

    def test_haversine_equator_distance(self, generator):
        """Test haversine distance along equator"""
        # 1 degree longitude at equator is approximately 111 km
        distance = generator._haversine(0, 0, 0, 1)

        # Should be around 111 km
        assert 110 < distance < 112

    def test_haversine_north_south(self, generator):
        """Test haversine distance north-south"""
        # 1 degree latitude is approximately 111 km
        distance = generator._haversine(0, 0, 1, 0)

        # Should be around 111 km
        assert 110 < distance < 112

    def test_sanitize_basic(self, generator):
        """Test path sanitization"""
        assert generator._sanitize("United States") == "united_states"
        assert generator._sanitize("New York") == "new_york"
        assert generator._sanitize("Los Angeles") == "los_angeles"

    def test_sanitize_special_chars(self, generator):
        """Test sanitization removes special characters"""
        assert generator._sanitize("São Paulo") == "so_paulo"
        assert generator._sanitize("Test/Path") == "testpath"
        assert generator._sanitize("Test\\Path") == "testpath"
        assert generator._sanitize("Test:Path") == "testpath"
        assert generator._sanitize("Test*Path") == "testpath"
        assert generator._sanitize("Test?Path") == "testpath"
        assert generator._sanitize("Test<Path>") == "testpath"

    def test_sanitize_multiple_spaces(self, generator):
        """Test sanitization with multiple spaces"""
        result = generator._sanitize("United  States  of  America")
        assert "_" in result
        assert result.islower()

    def test_basename_with_number(self, generator):
        """Test basename extraction from server name with number"""
        server = Server(
            name="us1234",
            hostname="us1234.nordvpn.com",
            station="1.2.3.4",
            load=50,
            country="United States",
            country_code="us",
            city="New York",
            latitude=40.7128,
            longitude=-74.0060,
            public_key="testkey123",
            distance=100.0
        )

        basename = generator._basename(server)
        assert basename.endswith(".conf")
        assert "1234" in basename

    def test_basename_without_number(self, generator):
        """Test basename extraction from server name without number"""
        server = Server(
            name="testserver",
            hostname="test.nordvpn.com",
            station="10.20.30.40",
            load=30,
            country="Japan",
            country_code="jp",
            city="Tokyo",
            latitude=35.6762,
            longitude=139.6503,
            public_key="testkey456",
            distance=200.0
        )

        basename = generator._basename(server)
        assert basename.endswith(".conf")
        # Should use station IP in fallback
        assert "wg" in basename


class TestGeneratorParseOne:
    """Test suite for _parse_one method"""

    @pytest.fixture
    def generator(self):
        """Create a generator instance for testing"""
        client = NordClient()
        ui = ConsoleManager()
        return Generator(client, ui)

    def test_parse_one_valid_server(self, generator):
        """Test parsing a valid server entry"""
        data = {
            "name": "us1234",
            "hostname": "us1234.nordvpn.com",
            "station": "192.168.1.1",
            "load": 50,
            "locations": [{
                "country": {
                    "name": "United States",
                    "code": "US",
                    "city": {"name": "New York"}
                },
                "latitude": 40.7128,
                "longitude": -74.0060
            }],
            "specifications": [{
                "identifier": "version",
                "values": [{"value": "2.1.0"}]
            }],
            "technologies": [{
                "identifier": "wireguard_udp",
                "metadata": [{
                    "name": "public_key",
                    "value": "testpublickey123="
                }]
            }]
        }

        server = generator._parse_one(data, 40.0, -73.0)

        assert server is not None
        assert server.name == "us1234"
        assert server.hostname == "us1234.nordvpn.com"
        assert server.station == "192.168.1.1"
        assert server.load == 50
        assert server.country == "United States"
        assert server.city == "New York"
        assert server.public_key == "testpublickey123="
        assert server.distance > 0

    def test_parse_one_missing_locations(self, generator):
        """Test parsing fails when locations are missing"""
        data = {
            "name": "us1234",
            "hostname": "us1234.nordvpn.com",
            "locations": [],
            "specifications": [{"identifier": "version", "values": [{"value": "2.1.0"}]}],
            "technologies": [{"identifier": "wireguard_udp", "metadata": [{"name": "public_key", "value": "key"}]}]
        }

        server = generator._parse_one(data, 40.0, -73.0)
        assert server is None

    def test_parse_one_old_version(self, generator):
        """Test parsing fails for old server versions"""
        data = {
            "name": "us1234",
            "hostname": "us1234.nordvpn.com",
            "station": "192.168.1.1",
            "locations": [{
                "country": {"name": "US", "code": "US", "city": {"name": "NY"}},
                "latitude": 40.0,
                "longitude": -73.0
            }],
            "specifications": [{
                "identifier": "version",
                "values": [{"value": "1.0.0"}]
            }],
            "technologies": [{
                "identifier": "wireguard_udp",
                "metadata": [{"name": "public_key", "value": "key"}]
            }]
        }

        server = generator._parse_one(data, 40.0, -73.0)
        assert server is None

    def test_parse_one_missing_public_key(self, generator):
        """Test parsing fails when public key is missing"""
        data = {
            "name": "us1234",
            "hostname": "us1234.nordvpn.com",
            "station": "192.168.1.1",
            "locations": [{
                "country": {"name": "US", "code": "US", "city": {"name": "NY"}},
                "latitude": 40.0,
                "longitude": -73.0
            }],
            "specifications": [{"identifier": "version", "values": [{"value": "2.1.0"}]}],
            "technologies": []
        }

        server = generator._parse_one(data, 40.0, -73.0)
        assert server is None

    def test_parse_one_malformed_data(self, generator):
        """Test parsing handles malformed data gracefully"""
        # Missing required fields
        data = {"name": "test"}

        server = generator._parse_one(data, 40.0, -73.0)
        assert server is None
