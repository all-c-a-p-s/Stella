use std::fs::File;
use std::io::Write;
use std::process::Command;
use std::{env, fs};

#[derive(Debug)]
struct Args {
    command: String,
    path: String,
    target: String,
}

impl Args {
    fn new(args: Vec<String>) -> Args {
        if args.len() != 5 {
            panic!("invalid number of arguments in command {:?}", args)
        }
        if &args[1] != "stella" {
            panic!("stella command not used")
        }
        let c = String::from(args[2].as_str());
        let p = String::from(args[3].as_str());
        let t = String::from(args[4].as_str());

        Args {
            command: c,
            path: p,
            target: t,
        }
    }
}

fn transpile(args: &Args) -> String {
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
    String::from_utf8(output.stdout).expect("output not valid UTF-8 string")
}

fn main() -> std::io::Result<()> {
    let args: Vec<String> = env::args().collect();
    let command_args = Args::new(args);

    let transpiled: String = transpile(&command_args);

    let already_exists = File::open(&command_args.target);
    if already_exists.is_ok() {
        //file already exists -> delete it and rewrite
        let delete = fs::remove_file(&command_args.target);
        match delete {
            Ok(()) => (),
            Err(_) => panic!(
                "failed to delete file {} before creating new file",
                &command_args.target
            ),
        }
    }

    let mut file = File::create(&command_args.target)
        .unwrap_or_else(|_| panic!("failed to create file {}", &command_args.target));

    file.write_all(transpiled.as_bytes())?;

    println!("Written transpiled to path {}", command_args.target);
    Ok(())
}
