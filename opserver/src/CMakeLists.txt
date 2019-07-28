cmake_minimum_required(VERSION 3.5.1)

project(opserver C CXX)

if(NOT MSVC)
  set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")
else()
  add_definitions(-D_WIN32_WINNT=0x600)
endif()

# This branch assumes that Eigen as well as gRPC and all their dependencies are
# already installed on this system, so they can be located by find_package().

# Find Protobuf installation
# Looks for protobuf-config.cmake file installed by Protobuf's cmake installation.
set(protobuf_MODULE_COMPATIBLE TRUE)
find_package(Protobuf CONFIG REQUIRED)
message(STATUS "Using protobuf ${protobuf_VERSION}")

set(_PROTOBUF_LIBPROTOBUF protobuf::libprotobuf)
set(_PROTOBUF_PROTOC $<TARGET_FILE:protobuf::protoc>)

# Find gRPC installation
# Looks for gRPCConfig.cmake file installed by gRPC's cmake installation.
find_package(gRPC CONFIG REQUIRED)
message(STATUS "Using gRPC ${gRPC_VERSION}")

set(_GRPC_GRPCPP_UNSECURE gRPC::grpc++_unsecure)
set(_GRPC_CPP_PLUGIN_EXECUTABLE $<TARGET_FILE:gRPC::grpc_cpp_plugin>)

# Find Eigen installation
# Looks for Eigen3Config.cmake file installed by Eigen's cmake installation.
find_package(Eigen3 3.3 REQUIRED NO_MODULE)
set(_EIGEN3_EIGEN Eigen3::Eigen)
message(STATUS "Using Eigen3 ${Eigen3_VERSION}")

# Proto file
get_filename_component(ops_proto "proto/matmult.proto" ABSOLUTE)
get_filename_component(ops_proto_path "${ops_proto}" PATH)

# Generated sources
set(ops_proto_srcs "${CMAKE_CURRENT_BINARY_DIR}/matmult.pb.cc")
set(ops_proto_hdrs "${CMAKE_CURRENT_BINARY_DIR}/matmult.pb.h")
set(ops_grpc_srcs "${CMAKE_CURRENT_BINARY_DIR}/matmult.grpc.pb.cc")
set(ops_grpc_hdrs "${CMAKE_CURRENT_BINARY_DIR}/matmult.grpc.pb.h")
add_custom_command(
      OUTPUT "${ops_proto_srcs}" "${ops_proto_hdrs}" "${ops_grpc_srcs}" "${ops_grpc_hdrs}"
      COMMAND ${_PROTOBUF_PROTOC}
      ARGS --grpc_out "${CMAKE_CURRENT_BINARY_DIR}"
        --cpp_out "${CMAKE_CURRENT_BINARY_DIR}"
        -I "${ops_proto_path}"
        --plugin=protoc-gen-grpc="${_GRPC_CPP_PLUGIN_EXECUTABLE}"
        "${ops_proto}"
      DEPENDS "${ops_proto}")

# Include generated *.pb.h files
include_directories(SYSTEM "${CMAKE_CURRENT_BINARY_DIR}")

set(ops_add_srcs "${CMAKE_CURRENT_SOURCE_DIR}/perf_srvc_impl.cc")

# Finally - assemble the server
foreach(_target server)
  add_executable(${_target} "${_target}.cc"
    ${ops_proto_srcs}
    ${ops_grpc_srcs}
    ${ops_add_srcs})
  target_link_libraries(${_target}
    ${_GRPC_GRPCPP_UNSECURE}
    ${_PROTOBUF_LIBPROTOBUF}
    ${_EIGEN3_EIGEN})
  target_compile_options(${_target} PRIVATE
    $<$<OR:$<CXX_COMPILER_ID:Clang>,$<CXX_COMPILER_ID:AppleClang>,$<CXX_COMPILER_ID:GNU>>:
      -Werror
      -pedantic-errors
      -Wall
      -Wextra
      -Wpedantic
      # -Wconversion
      # -Wsign-conversion
      -Wno-error=unused-parameter>
    $<$<CXX_COMPILER_ID:MSVC>:
      /WX
      /W4>)
endforeach()
