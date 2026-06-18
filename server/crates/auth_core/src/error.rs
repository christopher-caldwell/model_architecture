use thiserror::Error;

#[derive(Debug, Error, PartialEq, Eq)]
pub enum AuthError {
    #[error("token expired")]
    ExpiredToken,

    #[error("invalid token")]
    InvalidToken,
}
