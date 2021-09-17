from typing import Dict

import pytest


@pytest.fixture
def context() -> Dict[str, str]:
    """
    Base context for generation a new project

    Returns:
        Dict[str, str]: base context to create a new project
    """
    return {
        "project_name": "HeyDude Test Project",
        "project_slug": "heydude_test_project",
        "project_description": "A test project description",
        "app_name": "heyDude",
        "module_name": "github.com/user/repo",
        "golangci_version": "1.40.0",
    }
