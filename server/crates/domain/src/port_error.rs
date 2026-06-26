type PortSource = Box<dyn std::error::Error + Send + Sync>;

#[derive(Debug, thiserror::Error)]
pub enum PortError {
    #[error("repository operation failed")]
    Repository {
        #[source]
        source: PortSource,
    },

    #[error("unit of work failed")]
    UnitOfWork {
        #[source]
        source: PortSource,
    },
}

impl PortError {
    pub fn repository(source: impl Into<PortSource>) -> Self {
        Self::Repository {
            source: source.into(),
        }
    }

    pub fn unit_of_work(source: impl Into<PortSource>) -> Self {
        Self::UnitOfWork {
            source: source.into(),
        }
    }
}

pub type PortResult<T> = Result<T, PortError>;
