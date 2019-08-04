#include "perf_srvc_impl.h"

#include <Eigen/Dense>

using Eigen::Index;
using Eigen::MatrixXd;

Status PerformerServiceImpl::MultiplyMatrices(ServerContext* context,
                                              const PerformerRequest* request,
                                              PerformerResponse* response) {
  // TODO(amwolff): could be done in one loop with more logic.
  const auto& mplier = request->multiplier();
  const int mplier_columns_size = mplier.columns_size();
  const int mplier_rows_size = mplier.columns(0).coefficients_size();
  MatrixXd mplier_eigen(mplier_rows_size, mplier_columns_size);
  for (int i = 0; i != mplier_columns_size; ++i) {
    for (int j = 0; j != mplier_rows_size; ++j) {
      mplier_eigen(j, i) = mplier.columns(i).coefficients(j);
    }
  }
  const auto& mplicand = request->multiplicand();
  const int mplicand_columns_size = mplicand.columns_size();
  const int mplicand_rows_size = mplicand.columns(0).coefficients_size();
  MatrixXd mplicand_eigen(mplicand_rows_size, mplicand_columns_size);
  for (int i = 0; i != mplicand_columns_size; ++i) {
    for (int j = 0; j != mplicand_rows_size; ++j) {
      mplicand_eigen(j, i) = mplicand.columns(i).coefficients(j);
    }
  }

  MatrixXd result_eigen = mplier_eigen * mplicand_eigen;

  for (Index ic = 0, ec = result_eigen.cols(); ec != ic; ++ic) {
    auto response_result_column = response->mutable_result()->add_columns();
    for (Index ir = 0, er = result_eigen.rows(); er != ir; ++ir) {
      response_result_column->add_coefficients(result_eigen.coeff(ir, ic));
    }
  }

  return Status::OK;
}
