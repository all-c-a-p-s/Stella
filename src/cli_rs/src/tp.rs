use std::fs;
use std::fs::File;
use std::process::Command;

use std::env;
use std::io::Write;

use crate::error_parser::parse_error;
use crate::Args;

const BACKSLASH_ASCII: char = 98u8 as char;
//also string escape character lol

pub fn get_diretory(path: &String) -> String {
    let current_path = match env::current_dir() {
        Ok(path) => path,
        Err(_) => panic!("no directory found and failed to get working directory"),
    };

    if path.is_empty() {
        panic!("path {} has no filename", &path)
    }

    let mut directory_end: usize = path.len() - 1;

    for c in path.chars().rev() {
        match c {
            '/' | BACKSLASH_ASCII => break,
            _ => directory_end -= 1,
        };
    }

    if directory_end != 0 {
        return path[..directory_end].to_string();
    }

    format!("{}", current_path.display())
}

pub fn transpile(args: &Args) -> Result<String, String> {
    if args.command != "tp" {
        panic!("invalid command")
    }

    //TODO: write path to a .txt file
    //read file in main.go
    //then delete file again

    let current_directory: String = match env::current_dir() {
        Ok(path) => format!("{}", path.display()),
        Err(_) => panic!("failed to get current working directory"),
    };

    let compiler_directory = String::from("C:/Users/vajol/Documents/goProjects/Stella/src/cli");
    let ok = env::set_current_dir(&compiler_directory);
    if ok.is_err() {
        panic!("error entering directory {}", &compiler_directory)
    }

    let metadata = current_directory.clone() + "/" + args.path.as_str();

    match write_transpiled(metadata, String::from("metadata.txt")) {
        Ok(file) => file,
        Err(_) => panic!("error creating file metadata.txt"),
    };

    let ok = env::set_current_dir(&compiler_directory);
    if ok.is_err() {
        panic!("error entering directory {}", &compiler_directory)
    }

    let output = if cfg!(target_os = "windows") {
        Command::new("cmd")
            .args(["/C", "cli.exe"])
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
        println!("{}", String::from_utf8(output.stderr.clone()).unwrap());
        let msg: String = String::from_utf8(output.stderr).expect("failed to get error message");
        return Err(msg);
    }

    let ok = env::set_current_dir(&current_directory);
    if ok.is_err() {
        panic!("error entering directory {}", &current_directory)
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
