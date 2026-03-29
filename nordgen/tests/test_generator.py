"""Tests for generator.py - pure helper functions."""
import sys
import os
import math
import pytest
from unittest.mock import MagicMock, patch, AsyncMock

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))

from nord_config_generator.generator import Generator
from nord_config_generator.models import Server, UserPreferences


def _make_generator():
    client = MagicMock()
    ui = MagicMock()
    ui.add_task = MagicMock(return_value=1)
    ui.update_progress = MagicMock()
    ui.start_progress = MagicMock()
    ui.stop_progress = MagicMock()
    return Generator(client, ui)


# ---------------------------------------------------------------------------
# _check_version
# ---------------------------------------------------------------------------
class TestCheckVersion:
    def setup_method(self):
        self.gen = _make_generator()

    def test_exact_min_version(self):
        assert self.gen._check_version('2.1.0') is True

    def test_version_above_min_minor(self):
        assert self.gen._check_version('2.2.0') is True

    def test_version_below_min(self):
        assert self.gen._check_version('2.0.0') is False

    def test_version_1x(self):
        assert self.gen._check_version('1.9.9') is False

    def test_version_major_3(self):
        assert self.gen._check_version('3.0.0') is True

    def test_version_major_10(self):
        assert self.gen._check_version('10.0.0') is True

    def test_too_short_string(self):
        assert self.gen._check_version('2') is False

    def test_empty_string(self):
        assert self.gen._check_version('') is False

    def test_invalid_format_no_dot(self):
        assert self.gen._check_version('210') is False

    def test_non_numeric(self):
        assert self.gen._check_version('a.b.c') is False

    def test_two_part_version(self):
        assert self.gen._check_version('2.1') is True

    def test_two_part_version_below(self):
        assert self.gen._check_version('2.0') is False


# ---------------------------------------------------------------------------
# _haversine
# ---------------------------------------------------------------------------
class TestHaversine:
    def setup_method(self):
        self.gen = _make_generator()

    def test_same_point_is_zero(self):
        d = self.gen._haversine(0, 0, 0, 0)
        assert d == pytest.approx(0.0, abs=1e-6)

    def test_known_distance_ny_london(self):
        # New York (40.7128, -74.0060) to London (51.5074, -0.1278)
        d = self.gen._haversine(40.7128, -74.0060, 51.5074, -0.1278)
        assert d == pytest.approx(5570, rel=0.01)

    def test_antipodal_points(self):
        d = self.gen._haversine(0, 0, 0, 180)
        assert d == pytest.approx(math.pi * 6371, rel=0.001)

    def test_symmetry(self):
        d1 = self.gen._haversine(10.0, 20.0, 30.0, 40.0)
        d2 = self.gen._haversine(30.0, 40.0, 10.0, 20.0)
        assert d1 == pytest.approx(d2, rel=1e-10)

    def test_returns_float(self):
        d = self.gen._haversine(0, 0, 1, 1)
        assert isinstance(d, float)

    def test_positive_result(self):
        d = self.gen._haversine(0, 0, 10, 10)
        assert d > 0


# ---------------------------------------------------------------------------
# _sanitize
# ---------------------------------------------------------------------------
class TestSanitize:
    def setup_method(self):
        self.gen = _make_generator()

    def test_lowercase(self):
        assert self.gen._sanitize('United States') == 'united_states'

    def test_spaces_become_underscores(self):
        assert self.gen._sanitize('New York') == 'new_york'

    def test_strips_forbidden_chars(self):
        for ch in '<>:"/\\|?*':
            result = self.gen._sanitize(f'foo{ch}bar')
            assert ch not in result
            assert 'foobar' in result

    def test_null_byte_removed(self):
        result = self.gen._sanitize('foo\x00bar')
        assert '\x00' not in result

    def test_already_clean(self):
        assert self.gen._sanitize('us') == 'us'

    def test_empty_string(self):
        assert self.gen._sanitize('') == ''

    def test_mixed(self):
        result = self.gen._sanitize('Côte d\'Ivoire')
        assert result == result.lower()


# ---------------------------------------------------------------------------
# _basename
# ---------------------------------------------------------------------------
class TestBasename:
    def setup_method(self):
        self.gen = _make_generator()

    def _server(self, name, country_code='us', station='1.2.3.4'):
        return Server(
            name=name, hostname='host', station=station,
            load=0, country='Country', country_code=country_code,
            city='City', latitude=0, longitude=0,
            public_key='pk', distance=0
        )

    def test_standard_server_name(self):
        s = self._server('United States #1234', 'us')
        result = self.gen._basename(s)
        assert result == 'us1234.conf'

    def test_name_with_single_digit(self):
        s = self._server('US #1', 'us')
        result = self.gen._basename(s)
        assert result == 'us1.conf'

    def test_no_digits_falls_back_to_station(self):
        s = self._server('NoNumbers', 'us', '1.2.3.4')
        result = self.gen._basename(s)
        assert result.endswith('.conf')
        assert '.' not in result[:-5]

    def test_result_ends_with_conf(self):
        s = self._server('DE #999', 'de')
        result = self.gen._basename(s)
        assert result.endswith('.conf')

    def test_max_length_15_chars(self):
        s = self._server('US #123456789012', 'us')
        result = self.gen._basename(s)
        # basename without .conf should be ≤15 chars
        assert len(result) <= 20  # 15 + '.conf'

    def test_different_country_codes(self):
        for code in ['gb', 'de', 'jp', 'au']:
            s = self._server(f'Server #42', code)
            result = self.gen._basename(s)
            assert result.startswith(code)

    def test_fallback_truncated_to_15(self):
        # Station with many digits
        s = self._server('NoNum', 'us', '192.168.100.200')
        result = self.gen._basename(s)
        base = result[:-5]  # remove .conf
        assert len(base) <= 15


# ---------------------------------------------------------------------------
# _parse_one
# ---------------------------------------------------------------------------
def _make_raw_server(**overrides):
    data = {
        'name': 'US#1',
        'hostname': 'us1.nordvpn.com',
        'station': '1.2.3.4',
        'load': 25,
        'locations': [{
            'latitude': 40.7,
            'longitude': -74.0,
            'country': {
                'name': 'United States',
                'code': 'US',
                'city': {'name': 'New York'}
            }
        }],
        'technologies': [{
            'identifier': 'wireguard_udp',
            'metadata': [{'name': 'public_key', 'value': 'abc123pk='}]
        }],
        'specifications': [{
            'identifier': 'version',
            'values': [{'value': '2.1.0'}]
        }]
    }
    data.update(overrides)
    return data


class TestParseOne:
    def setup_method(self):
        self.gen = _make_generator()

    def test_valid_server(self):
        raw = _make_raw_server()
        s = self.gen._parse_one(raw, 0.0, 0.0)
        assert s is not None
        assert s.name == 'US#1'
        assert s.hostname == 'us1.nordvpn.com'
        assert s.station == '1.2.3.4'
        assert s.load == 25
        assert s.country == 'United States'
        assert s.country_code == 'us'
        assert s.city == 'New York'
        assert s.public_key == 'abc123pk='

    def test_missing_locations_returns_none(self):
        raw = _make_raw_server(locations=[])
        assert self.gen._parse_one(raw, 0.0, 0.0) is None

    def test_no_locations_key_returns_none(self):
        raw = _make_raw_server()
        del raw['locations']
        assert self.gen._parse_one(raw, 0.0, 0.0) is None

    def test_old_version_returns_none(self):
        raw = _make_raw_server(specifications=[{
            'identifier': 'version',
            'values': [{'value': '2.0.0'}]
        }])
        assert self.gen._parse_one(raw, 0.0, 0.0) is None

    def test_missing_public_key_returns_none(self):
        raw = _make_raw_server(technologies=[{
            'identifier': 'wireguard_udp',
            'metadata': []
        }])
        assert self.gen._parse_one(raw, 0.0, 0.0) is None

    def test_no_wireguard_tech_returns_none(self):
        raw = _make_raw_server(technologies=[{
            'identifier': 'openvpn_tcp',
            'metadata': [{'name': 'public_key', 'value': 'abc'}]
        }])
        assert self.gen._parse_one(raw, 0.0, 0.0) is None

    def test_country_code_lowercased(self):
        raw = _make_raw_server()
        s = self.gen._parse_one(raw, 0.0, 0.0)
        assert s.country_code == 'us'

    def test_distance_computed(self):
        raw = _make_raw_server()
        s = self.gen._parse_one(raw, 0.0, 0.0)
        assert isinstance(s.distance, float)
        assert s.distance >= 0

    def test_default_version_used_when_specs_empty(self):
        raw = _make_raw_server(specifications=[])
        # default version is "0.0.0" which fails check
        assert self.gen._parse_one(raw, 0.0, 0.0) is None

    def test_missing_load_defaults_to_zero(self):
        raw = _make_raw_server()
        del raw['load']
        s = self.gen._parse_one(raw, 0.0, 0.0)
        assert s is not None
        assert s.load == 0

    def test_corrupt_data_returns_none(self):
        assert self.gen._parse_one({}, 0.0, 0.0) is None


# ---------------------------------------------------------------------------
# _parse_batch
# ---------------------------------------------------------------------------
class TestParseBatch:
    def setup_method(self):
        self.gen = _make_generator()

    def test_empty_input(self):
        assert self.gen._parse_batch([], 0.0, 0.0) == []

    def test_valid_servers_returned(self):
        raw = [_make_raw_server(), _make_raw_server(name='US#2', station='2.2.2.2')]
        result = self.gen._parse_batch(raw, 0.0, 0.0)
        assert len(result) == 2

    def test_invalid_servers_skipped(self):
        raw = [_make_raw_server(), {'invalid': True}, _make_raw_server(locations=[])]
        result = self.gen._parse_batch(raw, 0.0, 0.0)
        assert len(result) == 1

    def test_all_invalid_returns_empty(self):
        raw = [{'bad': True}, {'also': 'bad'}]
        result = self.gen._parse_batch(raw, 0.0, 0.0)
        assert result == []
