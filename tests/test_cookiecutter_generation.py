import os
from typing import Any, Dict, List

import pytest
from hooks.post_gen_project import GRPC_FILES, OPENAPI_FILES
from pytest_cookies.plugin import Result

GRPC_SUPPORTED_OPTIONS: List[Dict[str, str]] = [
    {"grpc_enabled": "no"},
    {"grpc_enabled": "yes"},
]


@pytest.mark.parametrize("context_grpc", GRPC_SUPPORTED_OPTIONS)
def test_project_generation(
    cookies: Any,
    context: Dict[str, str],
    context_grpc: Dict[str, str],
) -> None:
    """
    Test that project is generated and fully rendered.

    Args:
        cookies (Any): cookiecutter generator test
        context (Dict[str, str]): base context for every project
        context_grpc (Dict[str, str]): grpc options
    """
    # generates a new project from your template based on the default values specified in cookiecutter.json
    result: Result = cookies.bake(
        extra_context={**context, **context_grpc}
    )

    # check generated project output
    assert result.exit_code == 0
    assert result.exception is None
    assert result.project.basename == context["project_slug"]
    assert result.project.isdir()

    # declare multiple set vars
    path: str

    ########
    # grpc #
    ########

    # check if all files for grpc was removed
    if result.context["grpc_enabled"] == "no":
        # check GRPC files
        for file_dir in GRPC_FILES:
            path = os.path.join(result.project, file_dir)
            assert os.path.exists(path) is False

        # check OpenAPI files
        for file_dir in OPENAPI_FILES:
            path = os.path.join(result.project, file_dir)
            assert os.path.exists(path) is True

    # check if all grpc files still exists
    if result.context["grpc_enabled"] == "yes":
        # check GRPC files
        for file_dir in GRPC_FILES:
            path = os.path.join(result.project, file_dir)
            assert os.path.exists(path) is True

        # check OpenAPI files
        for file_dir in OPENAPI_FILES:
            path = os.path.join(result.project, file_dir)
            assert os.path.exists(path) is False

    ######################
    # cookiecutter stuff #
    ######################

    # check if all template vars was rendert
    for root, _, files in os.walk(result.project, topdown=False):
        for file_name in files:
            file_path: str = os.path.join(root, file_name)
            with open(file_path) as file_content:
                is_template_tag: bool = "cookiecutter." in file_content.read()
                assert (
                    is_template_tag is False
                ), f"tempate var cookiecutter was not rendert in {file_path}"
