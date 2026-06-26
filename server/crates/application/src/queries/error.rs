use domain::PortError;

#[derive(Debug, thiserror::Error)]
pub enum QueryError {
    #[error(transparent)]
    Infrastructure(#[from] PortError),
}
