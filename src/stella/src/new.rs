use std::env;
use std::fs;
use std::fs::File;
use std::path::Path;

use crate::tp;
use crate::Args;

pub fn write_gomod(module_name: &str) -> String {
    String::from("module ") + module_name + "\n" + "\n" + "go 1.21.1"
}

pub fn create_module(module_name: &str) -> Result<(), String> {
    //create parent directory with same name as path argument in command
    //then create subdirectories src and tp
    let path = Path::new(module_name);
    if path.is_dir() {
        return Err(String::from("already exists"));
    }
    let create_ok = fs::create_dir(path);
    if create_ok.is_err() {
        eprintln!("failed to create directory {}", module_name);
        std::process::exit(1)
    }
    Ok(())
}

pub fn create_directory(dir_name: &str) -> std::io::Result<()> {
    let path = Path::new(dir_name);
    if path.is_dir() {
        eprintln!("src directory already exists");
        std::process::exit(1)
    }
    match fs::create_dir(path) {
        Ok(file) => Ok(file),
        Err(error) => panic!("error creating directory {}", error),
    }
}

pub fn create_file(filename: &str) -> std::io::Result<()> {
    let already_exists = File::open(filename);
    if already_exists.is_ok() {
        //file already exists -> delete it and rewrite
        let delete = fs::remove_file(filename);
        match delete {
            Ok(()) => (),
            Err(_) => panic!(
                "failed to delete file {} before creating new file",
                &filename
            ),
        }
    }

    match File::create(filename) {
        Ok(file) => file,
        Err(error) => panic!("failed to create file: {}", error),
    };

    Ok(())
}

pub fn create_subdirectories(module_name: &str) -> std::io::Result<()> {
    let ok = env::set_current_dir(module_name);
    if ok.is_err() {
        eprintln!("error entering directory {}", &module_name);
        std::process::exit(1)
    }
    create_directory("src")?;
    create_directory("tp")?;

    create_stella_files("src")?;

    let ok = env::set_current_dir("./..");
    if ok.is_err() {
        eprintln!("error exiting src directory");
        std::process::exit(1)
    }
    create_go_files("tp", module_name)?;

    Ok(())
}

pub fn create_stella_files(path: &str) -> std::io::Result<()> {
    let ok = env::set_current_dir(path);
    if ok.is_err() {
        eprintln!("error entering directory {}", &path);
        std::process::exit(1)
    }

    create_file("main.ste")?;
    Ok(())
}

pub fn create_go_files(path: &str, module_name: &str) -> std::io::Result<()> {
    let ok = env::set_current_dir(path);
    if ok.is_err() {
        eprintln!("error entering directory {}", &path);
        std::process::exit(1)
    }

    create_file("main.go")?;
    let gomod = write_gomod(module_name);
    tp::write_text(gomod, String::from("go.mod"))?;
    Ok(())
}

pub fn new(args: &Args) -> std::io::Result<()> {
    if args.command != "new" {
        panic!("new() called without new command")
    }
    if args.target.is_some() {
        eprintln!("new command used with target argument");
        std::process::exit(1)
    }
    match create_module(args.path.as_str()) {
        Ok(dir) => dir,
        Err(_) => panic!("error creating directory"),
    };

    create_subdirectories(args.path.as_str())?;
    println!("successfully created module {}", args.path.as_str());
    //this cd's into module, creates subdirectories and files

    Ok(())
}
