cmake_minimum_required(VERSION 3.5.1)

project(server C CXX)

if(NOT MSVC)
  set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")
else()
  add_definitions(-D_WIN32_WINNT=0x600)
endif()

find_package(Threads REQUIRED)

# This assumes that Eigen as well as gRPC and all their dependencies are already
# installed on this system, so they can be located by find_package().

# Find Protobuf installation
# Looks for protobuf-config.cmake file installed by Protobuf's cmake
# installation.
set(protobuf_MODULE_COMPATIBLE TRUE)
find_package(Protobuf CONFIG REQUIRED)
message(STATUS "Using protobuf ${protobuf_VERSION}")

set(_PROTOBUF_LIBPROTOBUF protobuf::libprotobuf)
set(_REFLECTION gRPC::grpc++_reflection)

if(CMAKE_CROSSCOMPILING)
  find_program(_PROTOBUF_PROTOC protoc)
else()
  set(_PROTOBUF_PROTOC $<TARGET_FILE:protobuf::protoc>)
endif()

# Find gRPC installation
# Looks for gRPCConfig.cmake file installed by gRPC's cmake installation.
find_package(gRPC CONFIG REQUIRED)
message(STATUS "Using gRPC ${gRPC_VERSION}")

set(_GRPC_GRPCPP gRPC::grpc++)

if(CMAKE_CROSSCOMPILING)
  find_program(_GRPC_CPP_PLUGIN_EXECUTABLE grpc_cpp_plugin)
else()
  set(_GRPC_CPP_PLUGIN_EXECUTABLE $<TARGET_FILE:gRPC::grpc_cpp_plugin>)
endif()

# Find Eigen installation
# Looks for Eigen3Config.cmake file installed by Eigen's cmake installation.
find_package(Eigen3 3.3 REQUIRED NO_MODULE)
message(STATUS "Using Eigen3 ${Eigen3_VERSION}")

set(_EIGEN3_EIGEN Eigen3::Eigen)

# Proto files
get_filename_component(api_proto "../../proto/amwolff/matmult/v1/matrix_product_api.proto" ABSOLUTE)
get_filename_component(matrix_proto "../../proto/amwolff/matmult/v1/matrix.proto" ABSOLUTE)
get_filename_component(proto_path "../../proto" ABSOLUTE)

# Generated sources
set(api_proto_srcs "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp/amwolff/matmult/v1/matrix_product_api.pb.cc")
set(api_proto_hdrs "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp/amwolff/matmult/v1/matrix_product_api.pb.h")
set(api_grpc_srcs "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp/amwolff/matmult/v1/matrix_product_api.grpc.pb.cc")
set(api_grpc_hdrs "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp/amwolff/matmult/v1/matrix_product_api.grpc.pb.h")
set(matrix_proto_srcs "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp/amwolff/matmult/v1/matrix.pb.cc")
set(matrix_proto_hdrs "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp/amwolff/matmult/v1/matrix.pb.h")
set(matrix_grpc_srcs "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp/amwolff/matmult/v1/matrix.grpc.pb.cc")
set(matrix_grpc_hdrs "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp/amwolff/matmult/v1/matrix.grpc.pb.h")
add_custom_command(
      OUTPUT "${api_proto_srcs}" "${api_proto_hdrs}" "${api_grpc_srcs}" "${api_grpc_hdrs}" "${matrix_proto_srcs}" "${matrix_proto_hdrs}" "${matrix_grpc_srcs}" "${matrix_grpc_hdrs}"
      COMMAND ${_PROTOBUF_PROTOC}
      ARGS --grpc_out "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp"
        --cpp_out "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp"
        -I "${proto_path}"
        --plugin=protoc-gen-grpc="${_GRPC_CPP_PLUGIN_EXECUTABLE}"
        "${api_proto}" "${matrix_proto}"
      DEPENDS "${api_proto}" "${matrix_proto}")

# Include generated *.pb.h files
include_directories(SYSTEM "${CMAKE_CURRENT_SOURCE_DIR}/../../gen/cpp")

set(srv_srcs "${CMAKE_CURRENT_SOURCE_DIR}/matmult_impl.cc")

# Finally - assemble the server
foreach(_target server)
  add_executable(${_target} "${_target}.cc"
    ${api_proto_srcs}
    ${api_grpc_srcs}
    ${matrix_proto_srcs}
    ${matrix_grpc_srcs}
    ${srv_srcs})
  target_link_libraries(${_target}
    ${_REFLECTION}
    ${_GRPC_GRPCPP}
    ${_PROTOBUF_LIBPROTOBUF}
    ${_EIGEN3_EIGEN}
    $<$<CXX_COMPILER_ID:GNU>:-fopenmp>)
  target_compile_options(${_target} PRIVATE
    $<$<OR:$<CXX_COMPILER_ID:Clang>,$<CXX_COMPILER_ID:AppleClang>,$<CXX_COMPILER_ID:GNU>>:
      -Werror
      -pedantic-errors
      -Wall
      -Wextra
      -Wconversion
      -Wsign-conversion
      -Wpedantic
      -Wno-error=conversion
      -Wno-error=sign-conversion
      -Wno-error=unused-parameter
      -march=native
      $<$<CXX_COMPILER_ID:GNU>:-fopenmp>>
    $<$<CXX_COMPILER_ID:MSVC>: # FIXME(amwolff): I have never tested this.
      # https://stackoverflow.com/questions/30800242/equivalents-to-gcc-clangs-march-native-in-other-compilers
      /WX
      /W4
      /openmp>)
endforeach()
