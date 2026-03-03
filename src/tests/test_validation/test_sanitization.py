"""Tests for input sanitization functions."""

import pytest


class TestInputSanitization:
    """Test input sanitization functions."""
    
    def test_sanitize_string(self):
        """Test string sanitization."""
        from validation.sanitize import SanitizeString
        
        # Test HTML escaping
        assert SanitizeString("<script>alert('xss')</script>") == "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
        
        # Test whitespace trimming
        assert SanitizeString("  test  ") == "test"
    
    def test_sanitize_sql_identifier(self):
        """Test SQL identifier sanitization."""
        from validation.sanitize import SanitizeSQLIdentifier
        
        # Test removal of dangerous characters
        assert SanitizeSQLIdentifier("table_name; DROP TABLE") == "table_nameDROPTABLE"
        assert SanitizeSQLIdentifier("valid_table_123") == "valid_table_123"
    
    def test_sanitize_filename(self):
        """Test filename sanitization."""
        from validation.sanitize import SanitizeFilename
        
        # Test path traversal prevention
        assert "../etc/passwd" not in SanitizeFilename("../etc/passwd")
        assert "file.txt" in SanitizeFilename("file.txt")
    
    def test_sanitize_url(self):
        """Test URL sanitization."""
        from validation.sanitize import SanitizeURL
        
        # Test dangerous protocol removal
        assert SanitizeURL("javascript:alert('xss')") == ""
        assert SanitizeURL("https://example.com") == "https://example.com"
    
    def test_sanitize_email(self):
        """Test email sanitization."""
        from validation.sanitize import SanitizeEmail
        
        # Test whitespace removal and lowercasing
        assert SanitizeEmail("  Test@Example.COM  ") == "test@example.com"

