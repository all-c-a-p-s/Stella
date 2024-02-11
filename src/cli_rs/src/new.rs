use std::env;
use std::fs;
use std::fs::File;
use std::path::Path;

use crate::Args;

pub fn create_new(p: &str) -> Result<(), String> {
    let path = Path::new(p);
    if path.is_dir() {
        return Err(String::from("already exists"));
    }
    match fs::create_dir(path) {
        Ok(file) => Ok(file),
        Err(error) => panic!("error creating directory {}", error),
    }
}

pub fn create_files(path: &str) -> std::io::Result<()> {
    let ok = env::set_current_dir(path);
    if ok.is_err() {
        panic!("error entering directory {}", &path)
    }

    let already_exists = File::open("main.ste");
    if already_exists.is_ok() {
        //file already exists -> delete it and rewrite
        let delete = fs::remove_file("main.ste");
        match delete {
            Ok(()) => (),
            Err(_) => panic!("failed to delete file main.ste before creating new file"),
        }
    }

    match File::create("main.ste") {
        Ok(file) => file,
        Err(error) => panic!("failed to create file: {}", error),
    };

    let already_exists = File::open("stella.toml");
    if already_exists.is_ok() {
        //file already exists -> delete it and rewrite
        let delete = fs::remove_file("stella.toml");
        match delete {
            Ok(()) => (),
            Err(_) => panic!("failed to delete file stella.toml before creating new file"),
        }
    }

    match File::create("stella.toml") {
        Ok(file) => file,
        Err(error) => panic!("failed to create file: {}", error),
    };
    Ok(())
}

pub fn new(args: &Args) -> std::io::Result<()> {
    if args.command != "new" {
        panic!("new() called without new command")
    }
    if args.target.is_some() {
        panic!("new command used with target argument")
    }
    match create_new(args.path.as_str()) {
        Ok(dir) => dir,
        Err(_) => panic!("error creating directory"),
    };

    create_files(args.path.as_str())?;
    Ok(())
}
