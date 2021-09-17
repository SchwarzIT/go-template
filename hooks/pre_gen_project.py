import sys

if __name__ == "__main__":
    # remove grpc support, if selected
    if (
        "{{ cookiecutter.grpc_gateway_enabled }}" == "yes"
        and "{{ cookiecutter.grpc_enabled }}" == "no"
    ):
        print(
            "ERROR: grpc_enabled needs to be set to 'yes' to enable the grpc_gateway_enabled option!"
        )

        sys.exit(1)
