use std::fs;
use std::fs::File;
use std::process::Command;

use std::io::Write;

use crate::error_parser::parse_error;
use crate::Args;

pub fn transpile(args: &Args) -> Result<String, String> {
    if args.command != "tp" {
        panic!("invalid command")
    }
    let output = if cfg!(target_os = "windows") {
        Command::new("../cli/cli.exe")
            .args([&args.path])
            .output()
            .expect("failed to execute process")
    } else {
        Command::new("sh")
            .current_dir("./../cli/")
            .arg("-c")
            .arg("go run .")
            .arg(&args.path)
            .output()
            .expect("failed to execute process")
    };
    if output.stdout.is_empty() {
        if output.stderr.is_empty() {
            panic!("both output.stdout and output.stderr empty")
        }
        let msg: String = String::from_utf8(output.stderr).expect("failed to get error message");
        return Err(msg);
    }
    Ok(String::from_utf8(output.stdout).expect("output not valid UTF-8 string"))
}

pub fn write_transpiled(tpd: String, target: String) -> std::io::Result<()> {
    let already_exists = File::open(&target);
    if already_exists.is_ok() {
        //file already exists -> delete it and rewrite
        let delete = fs::remove_file(&target);
        match delete {
            Ok(()) => (),
            Err(_) => panic!("failed to delete file {} before creating new file", target),
        }
    }

    let mut file = match File::create(&target) {
        Ok(file) => file,
        Err(error) => panic!("failed to create file: {}", error),
    };

    file.write_all(tpd.as_bytes())?;
    Ok(())
}

pub fn tp(args: &Args) -> std::io::Result<()> {
    if args.command != "tp" {
        panic!("tp() called without tp command")
    }

    let target = match args.target {
        Some(ref path) => path,
        None => panic!("tp command used with no target path"),
    };

    let transpiled = match transpile(args) {
        Ok(tp) => tp,
        Err(msg) => panic!("Transpiler Error: {:?}", parse_error(msg)),
    };

    write_transpiled(transpiled, target.to_owned())?;

    println!("Transpiled to path {}", target);

    Ok(())
}
