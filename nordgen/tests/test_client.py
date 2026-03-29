"""Tests for client.py - NordClient HTTP client (mocked)."""
import sys
import os
import pytest
from unittest.mock import AsyncMock, MagicMock, patch

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))

from nord_config_generator.client import NordClient


@pytest.fixture
async def client():
    async with NordClient() as c:
        yield c


# ---------------------------------------------------------------------------
# get_key
# ---------------------------------------------------------------------------
class TestGetKey:
    async def test_returns_key_on_success(self, client):
        mock_resp = AsyncMock()
        mock_resp.status = 200
        mock_resp.json = AsyncMock(return_value={'nordlynx_private_key': 'mykey123'})
        mock_resp.__aenter__ = AsyncMock(return_value=mock_resp)
        mock_resp.__aexit__ = AsyncMock(return_value=False)

        client._session.get = MagicMock(return_value=mock_resp)
        result = await client.get_key('a' * 64)
        assert result == 'mykey123'

    async def test_returns_none_on_non_200(self, client):
        mock_resp = AsyncMock()
        mock_resp.status = 401
        mock_resp.__aenter__ = AsyncMock(return_value=mock_resp)
        mock_resp.__aexit__ = AsyncMock(return_value=False)

        client._session.get = MagicMock(return_value=mock_resp)
        result = await client.get_key('a' * 64)
        assert result is None

    async def test_returns_none_when_no_session(self):
        c = NordClient()
        result = await c.get_key('a' * 64)
        assert result is None

    async def test_returns_none_on_network_error(self, client):
        import aiohttp
        client._session.get = MagicMock(side_effect=aiohttp.ClientError())
        result = await client.get_key('a' * 64)
        assert result is None

    async def test_auth_header_uses_base64(self, client):
        import base64
        captured = {}

        mock_resp = AsyncMock()
        mock_resp.status = 200
        mock_resp.json = AsyncMock(return_value={'nordlynx_private_key': 'key'})
        mock_resp.__aenter__ = AsyncMock(return_value=mock_resp)
        mock_resp.__aexit__ = AsyncMock(return_value=False)

        def capture_get(url, headers=None, **kwargs):
            captured['headers'] = headers
            return mock_resp

        client._session.get = capture_get
        await client.get_key('token123')

        auth = captured['headers']['Authorization']
        assert auth.startswith('Basic ')
        decoded = base64.b64decode(auth[6:]).decode()
        assert decoded == 'token:token123'


# ---------------------------------------------------------------------------
# get_servers
# ---------------------------------------------------------------------------
class TestGetServers:
    async def test_returns_list_on_success(self, client):
        servers = [{'name': 'US#1'}, {'name': 'DE#1'}]
        mock_resp = AsyncMock()
        mock_resp.status = 200
        mock_resp.json = AsyncMock(return_value=servers)
        mock_resp.__aenter__ = AsyncMock(return_value=mock_resp)
        mock_resp.__aexit__ = AsyncMock(return_value=False)

        client._session.get = MagicMock(return_value=mock_resp)
        result = await client.get_servers()
        assert result == servers

    async def test_returns_empty_on_non_200(self, client):
        mock_resp = AsyncMock()
        mock_resp.status = 500
        mock_resp.__aenter__ = AsyncMock(return_value=mock_resp)
        mock_resp.__aexit__ = AsyncMock(return_value=False)

        client._session.get = MagicMock(return_value=mock_resp)
        result = await client.get_servers()
        assert result == []

    async def test_returns_empty_when_no_session(self):
        c = NordClient()
        result = await c.get_servers()
        assert result == []

    async def test_returns_empty_on_network_error(self, client):
        import aiohttp
        client._session.get = MagicMock(side_effect=aiohttp.ClientError())
        result = await client.get_servers()
        assert result == []


# ---------------------------------------------------------------------------
# get_geo
# ---------------------------------------------------------------------------
class TestGetGeo:
    async def test_returns_lat_lon_on_success(self, client):
        mock_resp = AsyncMock()
        mock_resp.status = 200
        mock_resp.json = AsyncMock(return_value={'latitude': 40.7, 'longitude': -74.0})
        mock_resp.__aenter__ = AsyncMock(return_value=mock_resp)
        mock_resp.__aexit__ = AsyncMock(return_value=False)

        client._session.get = MagicMock(return_value=mock_resp)
        lat, lon = await client.get_geo()
        assert lat == pytest.approx(40.7)
        assert lon == pytest.approx(-74.0)

    async def test_returns_zeros_on_non_200(self, client):
        mock_resp = AsyncMock()
        mock_resp.status = 403
        mock_resp.__aenter__ = AsyncMock(return_value=mock_resp)
        mock_resp.__aexit__ = AsyncMock(return_value=False)

        client._session.get = MagicMock(return_value=mock_resp)
        lat, lon = await client.get_geo()
        assert lat == 0.0
        assert lon == 0.0

    async def test_returns_zeros_when_no_session(self):
        c = NordClient()
        lat, lon = await c.get_geo()
        assert lat == 0.0
        assert lon == 0.0

    async def test_returns_zeros_on_network_error(self, client):
        import aiohttp
        client._session.get = MagicMock(side_effect=aiohttp.ClientError())
        lat, lon = await client.get_geo()
        assert lat == 0.0
        assert lon == 0.0

    async def test_returns_zeros_on_missing_fields(self, client):
        mock_resp = AsyncMock()
        mock_resp.status = 200
        mock_resp.json = AsyncMock(return_value={})
        mock_resp.__aenter__ = AsyncMock(return_value=mock_resp)
        mock_resp.__aexit__ = AsyncMock(return_value=False)

        client._session.get = MagicMock(return_value=mock_resp)
        lat, lon = await client.get_geo()
        assert lat == 0.0
        assert lon == 0.0
