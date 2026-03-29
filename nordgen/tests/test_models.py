"""Tests for models.py - Server, UserPreferences, Stats data classes."""
import pytest
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))

from nord_config_generator.models import Server, UserPreferences, Stats


class TestServer:
    def _make(self, **kwargs):
        defaults = dict(
            name='US#1', hostname='us1.nordvpn.com', station='1.2.3.4',
            load=20, country='United States', country_code='us',
            city='New York', latitude=40.7128, longitude=-74.0060,
            public_key='abc123', distance=100.0
        )
        defaults.update(kwargs)
        return Server(**defaults)

    def test_basic_creation(self):
        s = self._make()
        assert s.name == 'US#1'
        assert s.hostname == 'us1.nordvpn.com'
        assert s.station == '1.2.3.4'
        assert s.load == 20
        assert s.country == 'United States'
        assert s.country_code == 'us'
        assert s.city == 'New York'
        assert s.latitude == 40.7128
        assert s.longitude == -74.0060
        assert s.public_key == 'abc123'
        assert s.distance == 100.0

    def test_slots_prevents_new_attributes(self):
        s = self._make()
        with pytest.raises(AttributeError):
            s.new_field = 'value'

    def test_zero_load(self):
        s = self._make(load=0)
        assert s.load == 0

    def test_max_load(self):
        s = self._make(load=100)
        assert s.load == 100

    def test_zero_distance(self):
        s = self._make(distance=0.0)
        assert s.distance == 0.0

    def test_large_distance(self):
        s = self._make(distance=20000.0)
        assert s.distance == 20000.0

    def test_special_chars_in_name(self):
        s = self._make(name='US #123')
        assert s.name == 'US #123'


class TestUserPreferences:
    def test_defaults(self):
        p = UserPreferences()
        assert p.dns == '103.86.96.100'
        assert p.use_ip is False
        assert p.keepalive == 25

    def test_custom_values(self):
        p = UserPreferences(dns='8.8.8.8', use_ip=True, keepalive=60)
        assert p.dns == '8.8.8.8'
        assert p.use_ip is True
        assert p.keepalive == 60

    def test_slots_prevents_new_attributes(self):
        p = UserPreferences()
        with pytest.raises(AttributeError):
            p.extra = 'x'

    def test_multiple_dns(self):
        p = UserPreferences(dns='8.8.8.8,8.8.4.4')
        assert p.dns == '8.8.8.8,8.8.4.4'

    def test_min_keepalive(self):
        p = UserPreferences(keepalive=15)
        assert p.keepalive == 15

    def test_max_keepalive(self):
        p = UserPreferences(keepalive=120)
        assert p.keepalive == 120


class TestStats:
    def test_initial_zeros(self):
        s = Stats()
        assert s.total == 0
        assert s.best == 0
        assert s.rejected == 0

    def test_assignment(self):
        s = Stats()
        s.total = 100
        s.best = 50
        s.rejected = 10
        assert s.total == 100
        assert s.best == 50
        assert s.rejected == 10

    def test_slots_prevents_new_attributes(self):
        s = Stats()
        with pytest.raises(AttributeError):
            s.extra = 'x'
