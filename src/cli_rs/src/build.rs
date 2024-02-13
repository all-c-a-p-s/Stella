use std::env;
use std::process::Command;

use crate::error_parser::parse_error;
use crate::Args;

pub fn go_build(path: &str) -> Result<String, String> {
    let ok = env::set_current_dir(path);
    if ok.is_err() {
        panic!("error entering module directory {}", &path)
    }

    let ok = env::set_current_dir("tp");
    if ok.is_err() {
        panic!("error entering tp directory {}", &path)
    }

    let output = if cfg!(target_os = "windows") {
        Command::new("cmd")
            .args(["/C", "go build ."])
            .output()
            .expect("failed to execute process")
    } else {
        Command::new("sh")
            .current_dir(path)
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
        panic!("build() called without build command")
    }

    if args.target.is_some() {
        panic!("stella build command used with unexpected target parameter")
    }

    let status = match go_build(args.path.as_str()) {
        Ok(tp) => tp,
        Err(msg) => panic!("Go Compilation Error: {:?}", parse_error(msg)),
    };

    println!("{}", status);

    Ok(())
}
