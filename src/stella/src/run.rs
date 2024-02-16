use std::env;
use std::fs::File;
use std::process::Command;

use crate::error_parser::parse_error;
use crate::Args;

pub fn go_run(module_name: &str) -> Result<String, String> {
    let ok = env::set_current_dir(module_name);
    if ok.is_err() {
        panic!("error entering module directory {}", &module_name)
    }

    let ok = env::set_current_dir("tp");
    if ok.is_err() {
        panic!("error entering tp directory {}", &module_name)
    }

    let exe_name = String::from(module_name) + ".exe";

    let already_exists = File::open(&exe_name);
    if already_exists.is_ok() {
        //file already exists -> no need to recompile

        let output = if cfg!(target_os = "windows") {
            Command::new("cmd")
                .args(["/C", exe_name.as_str()])
                .output()
                .expect("failed to execute process")
        } else {
            Command::new("sh")
                .arg(exe_name.as_str())
                .output()
                .expect("failed to execute process")
        };

        if output.stderr.is_empty() {
            let res: String =
                String::from_utf8(output.stdout).unwrap_or(String::from("error getting output"));
            return Ok(res);
        }
    }

    Err(format!("found no executable file in module {}. Try generating an executable with to stella build {} command", module_name, module_name))
}

pub fn run(args: &Args) -> std::io::Result<()> {
    if args.command != "run" {
        panic!("run() called without run command")
    }

    if args.target.is_some() {
        panic!("stella run command used with unexpected target parameter")
    }

    let status = match go_run(args.path.as_str()) {
        Ok(tp) => tp,
        Err(msg) => panic!("Go Compilation Error: {:?}", parse_error(msg)),
    };

    println!("{}", status);
    Ok(())
}
