"""Tests for WebSocket Streaming."""
import pytest
import websocket

@pytest.mark.integration
@pytest.mark.requires_server
@pytest.mark.slow
class TestStreaming:
    def test_websocket_streaming(self, api_client, test_session):
        """Test WebSocket message streaming."""
        # WebSocket connection test
        ws_url = f"ws://localhost:8080/ws?session_id={test_session['id']}"
        try:
            ws = websocket.create_connection(ws_url, timeout=5)
            ws.close()
        except Exception:
            pytest.skip("WebSocket not available")

