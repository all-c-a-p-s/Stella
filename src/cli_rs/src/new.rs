use std::fs;
use std::path::Path;

pub fn create_new(p: String) -> Result<(), String> {
    let path = Path::new(p.as_str());
    if path.is_dir() {
        return Err(String::from("already exists"));
    }
    match fs::create_dir(path) {
        Ok(file) => Ok(file),
        Err(error) => panic!("error creating directory {}", error),
    }
}
