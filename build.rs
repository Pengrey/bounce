use serde::Deserialize;
use std::collections::HashMap;
use std::env;
use std::fs;
use std::path::Path;

// Updated struct to match the new config.toml format
#[derive(Deserialize)]
struct Config {
    target_url: String,
    headers: HashMap<String, String>,
}

fn main() {
    println!("cargo:rerun-if-changed=config.toml");
    println!("cargo:rerun-if-changed=build.rs");

    let config_str = fs::read_to_string("config.toml").expect("Failed to read config.toml");
    let config: Config = toml::from_str(&config_str).expect("Failed to parse config.toml");

    let mut generated_code = String::new();

    // Generate constants for host and path
    generated_code.push_str(&format!(
        "pub const TARGET_URL: &'static str = \"{}\";\n",
        config.target_url
    ));

    // Generate the headers array
    generated_code.push_str("\npub const HEADERS: &'static [(&'static str, &'static str)] = &[\n");
    for (key, value) in config.headers {
        generated_code.push_str(&format!(
            "    (\"{}\", \"{}\"),\n",
                                         key.escape_default(),
                                         value.escape_default()
        ));
    }
    generated_code.push_str("];\n");

    let out_dir = env::var("OUT_DIR").unwrap();
    let dest_path = Path::new(&out_dir).join("config.rs");
    fs::write(&dest_path, generated_code).expect("Failed to write to config.rs");
}
