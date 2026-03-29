"""Tests for main.py - resolve_key function."""
import sys
import os
import pytest
from unittest.mock import MagicMock, AsyncMock

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))

from nord_config_generator.main import resolve_key


class TestResolveKey:
    def _make_ui_client(self, token_from_prompt='a' * 64, key_response='valid_key'):
        ui = MagicMock()
        ui.prompt_secret = MagicMock(return_value=token_from_prompt)
        ui.spin = MagicMock()
        ui.success = MagicMock()
        ui.fail = MagicMock()
        ui.error = MagicMock()

        client = MagicMock()
        client.get_key = AsyncMock(return_value=key_response)
        return ui, client

    async def test_valid_token_returns_key(self):
        ui, client = self._make_ui_client()
        token = 'a' * 64
        result = await resolve_key(ui, client, token)
        assert result == 'valid_key'
        client.get_key.assert_called_once_with(token)

    async def test_prompts_when_token_empty(self):
        ui, client = self._make_ui_client(token_from_prompt='b' * 64)
        result = await resolve_key(ui, client, '')
        ui.prompt_secret.assert_called_once()
        assert result == 'valid_key'

    async def test_invalid_token_length_returns_empty(self):
        ui, client = self._make_ui_client()
        result = await resolve_key(ui, client, 'short')
        assert result == ''
        ui.error.assert_called_once()
        client.get_key.assert_not_called()

    async def test_token_63_chars_returns_empty(self):
        ui, client = self._make_ui_client()
        result = await resolve_key(ui, client, 'a' * 63)
        assert result == ''

    async def test_api_returns_none_returns_empty(self):
        ui, client = self._make_ui_client(key_response=None)
        result = await resolve_key(ui, client, 'a' * 64)
        assert result == ''
        ui.fail.assert_called_once()

    async def test_token_65_chars_returns_empty(self):
        ui, client = self._make_ui_client()
        result = await resolve_key(ui, client, 'a' * 65)
        assert result == ''
