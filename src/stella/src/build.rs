use std::process::Command;
use std::{env, process};

use crate::error_parser::parse_error;
use crate::Args;

pub fn go_build(path: &str) -> Result<String, String> {
    let ok = env::set_current_dir(path);
    if ok.is_err() {
        eprintln!("error entering module directory {}", &path);
        process::exit(1)
    }

    let ok = env::set_current_dir("tp");
    if ok.is_err() {
        eprintln!("error entering tp directory {}", &path);
        process::exit(1)
    }

    let output = if cfg!(target_os = "windows") {
        Command::new("cmd")
            .args(["/C", "go build ."])
            .output()
            .expect("failed to execute process")
    } else {
        Command::new("sh")
            .arg("-c")
            .arg("go build .")
            .output()
            .expect("failed to execute process")
    };

    if output.stderr.is_empty() {
        return Ok(String::from("build successful"));
    }
    let msg: String = String::from_utf8(output.stderr).expect("failed to get error message");
    Err(msg)
}

pub fn build(args: &Args) -> std::io::Result<()> {
    if args.command != "build" {
        eprintln!("build() called without build command")
    }

    if args.target.is_some() {
        eprintln!("stella build command used with unexpected target parameter")
    }

    let status = match go_build(args.path.as_str()) {
        Ok(tp) => tp,
        Err(msg) => panic!("Go Compilation Error: {:?}", parse_error(msg)),
    };

    println!("{}", status);

    Ok(())
}
