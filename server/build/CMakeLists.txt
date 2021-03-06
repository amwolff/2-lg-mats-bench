cmake_minimum_required(VERSION 3.5.1)

project(server-SuperBuild C CXX)

include(ExternalProject)

# Builds c-ares project from the git submodule.
#
# Note: For all external projects, instead of using checked-out code, one could
# specify GIT_REPOSITORY and GIT_TAG to have cmake download the dependency
# directly, without needing to add a submodule to your project.
ExternalProject_Add(c-ares
  PREFIX c-ares
  SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../third_party/grpc/third_party/cares/cares"
  CMAKE_CACHE_ARGS
        -DCARES_SHARED:BOOL=OFF
        -DCARES_STATIC:BOOL=ON
        -DCARES_STATIC_PIC:BOOL=ON
        -DCMAKE_INSTALL_PREFIX:PATH=${CMAKE_CURRENT_BINARY_DIR}/c-ares
        -DCMAKE_BUILD_TYPE:STRING=$<CONFIG>
)

# Builds protobuf project from the git submodule.
ExternalProject_Add(protobuf
  PREFIX protobuf
  SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../third_party/grpc/third_party/protobuf/cmake"
  CMAKE_CACHE_ARGS
        -Dprotobuf_BUILD_TESTS:BOOL=OFF
        -Dprotobuf_WITH_ZLIB:BOOL=OFF
        -Dprotobuf_MSVC_STATIC_RUNTIME:BOOL=OFF
        -DCMAKE_INSTALL_PREFIX:PATH=${CMAKE_CURRENT_BINARY_DIR}/protobuf
        -DCMAKE_BUILD_TYPE:STRING=$<CONFIG>
)

# Builds zlib project from the git submodule.
ExternalProject_Add(zlib
  PREFIX zlib
  SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../third_party/grpc/third_party/zlib"
  CMAKE_CACHE_ARGS
        -DCMAKE_INSTALL_PREFIX:PATH=${CMAKE_CURRENT_BINARY_DIR}/zlib
        -DCMAKE_BUILD_TYPE:STRING=$<CONFIG>
)

# The location where protobuf-config.cmake will be installed varies by platform
if (WIN32)
  set(_FINDPACKAGE_PROTOBUF_CONFIG_DIR "${CMAKE_CURRENT_BINARY_DIR}/protobuf/cmake")
else()
  set(_FINDPACKAGE_PROTOBUF_CONFIG_DIR "${CMAKE_CURRENT_BINARY_DIR}/protobuf/lib/cmake/protobuf")
endif()

# if OPENSSL_ROOT_DIR is set, propagate that hint path to the external projects
# with OpenSSL dependency.
set(_CMAKE_ARGS_OPENSSL_ROOT_DIR "")
if (OPENSSL_ROOT_DIR)
  set(_CMAKE_ARGS_OPENSSL_ROOT_DIR "-DOPENSSL_ROOT_DIR:PATH=${OPENSSL_ROOT_DIR}")
endif()

# Builds gRPC project from the git submodule and set arguments so that all the
# dependencies are correctly located.
ExternalProject_Add(grpc
  PREFIX grpc
  SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../third_party/grpc"
  CMAKE_CACHE_ARGS
        -DgRPC_INSTALL:BOOL=ON
        -DgRPC_BUILD_TESTS:BOOL=OFF
        -DgRPC_PROTOBUF_PROVIDER:STRING=package
        -DgRPC_PROTOBUF_PACKAGE_TYPE:STRING=CONFIG
        -DProtobuf_DIR:PATH=${_FINDPACKAGE_PROTOBUF_CONFIG_DIR}
        -DgRPC_ZLIB_PROVIDER:STRING=package
        -DZLIB_ROOT:STRING=${CMAKE_CURRENT_BINARY_DIR}/zlib
        -DgRPC_CARES_PROVIDER:STRING=package
        -Dc-ares_DIR:PATH=${CMAKE_CURRENT_BINARY_DIR}/c-ares/lib/cmake/c-ares
        -DgRPC_SSL_PROVIDER:STRING=package
        ${_CMAKE_ARGS_OPENSSL_ROOT_DIR}
        -DCMAKE_INSTALL_PREFIX:PATH=${CMAKE_CURRENT_BINARY_DIR}/grpc
        -DCMAKE_BUILD_TYPE:STRING=$<CONFIG>
  DEPENDS c-ares protobuf zlib
)

# Builds eigen project from the git submodule.
ExternalProject_Add(eigen
  PREFIX eigen
  SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../third_party/eigen-git-mirror"
  CMAKE_CACHE_ARGS
        -DCMAKE_INSTALL_PREFIX:PATH=${CMAKE_CURRENT_BINARY_DIR}/eigen
        -DCMAKE_BUILD_TYPE:STRING=$<CONFIG>
)

# Build the server project itself using a CMakeLists.txt that assumes all the
# dependencies have already been installed.
#
# Even though server is not really an "external project" from perspective of
# this build, we are still importing it using ExternalProject_Add because that
# allows us to use find_package() to locate all the dependencies (if we were
# building server directly in this build, we would have needed to manually
# import the libraries as opposed to reusing targets exported by gRPC and
# protobuf).
ExternalProject_Add(server
  PREFIX server
  SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../src"
  BINARY_DIR "${CMAKE_CURRENT_BINARY_DIR}/server"
  INSTALL_COMMAND ""
  CMAKE_CACHE_ARGS
        -DProtobuf_DIR:PATH=${_FINDPACKAGE_PROTOBUF_CONFIG_DIR}
        -Dc-ares_DIR:PATH=${CMAKE_CURRENT_BINARY_DIR}/c-ares/lib/cmake/c-ares
        -DZLIB_ROOT:STRING=${CMAKE_CURRENT_BINARY_DIR}/zlib
        ${_CMAKE_ARGS_OPENSSL_ROOT_DIR}
        -DgRPC_DIR:PATH=${CMAKE_CURRENT_BINARY_DIR}/grpc/lib/cmake/grpc
        -DEigen3_DIR:PATH=${CMAKE_CURRENT_BINARY_DIR}/eigen/cmake
        -DCMAKE_BUILD_TYPE:STRING=$<CONFIG>
  DEPENDS protobuf grpc eigen
)
