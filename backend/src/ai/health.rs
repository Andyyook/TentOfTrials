use serde::{Deserialize, Serialize};
use std::time::{SystemTime, UNIX_EPOCH};

/// Structured health endpoint response with service metadata.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HealthResponse {
    pub status: String,
    pub version: String,
    pub commit: String,
    pub uptime_seconds: u64,
    pub features: Vec<String>,
    pub timestamp: String,
}

/// Collect health information from the backend subsystems.
pub fn health_status(start_time: &SystemTime) -> HealthResponse {
    let uptime = SystemTime::now()
        .duration_since(*start_time)
        .unwrap_or_default()
        .as_secs();

    let commit = std::env::var("GIT_COMMIT").unwrap_or_else(|_| "unknown".to_string());

    HealthResponse {
        status: "ok".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
        commit,
        uptime_seconds: uptime,
        features: vec![
            "discovery".to_string(),
            "messaging".to_string(),
            "registry".to_string(),
            "consensus".to_string(),
        ],
        timestamp: chrono::Utc::now().to_rfc3339(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::time::SystemTime;

    #[test]
    fn test_health_returns_ok_status() {
        let start = SystemTime::now();
        std::thread::sleep(std::time::Duration::from_millis(10));
        let health = health_status(&start);
        assert_eq!(health.status, "ok");
        assert!(health.uptime_seconds > 0);
    }

    #[test]
    fn test_health_contains_expected_fields() {
        let start = SystemTime::now();
        let health = health_status(&start);
        assert_eq!(health.version, env!("CARGO_PKG_VERSION"));
        assert!(!health.features.is_empty());
        assert!(health.features.contains(&"discovery".to_string()));
    }

    #[test]
    fn test_health_serializable() {
        let start = SystemTime::now();
        let health = health_status(&start);
        let json = serde_json::to_string(&health).unwrap();
        assert!(json.contains(""status":"ok""));
    }
}
