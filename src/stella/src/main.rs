use std::env;

pub mod build;
pub mod error_parser;
pub mod new;
pub mod run;
pub mod tp;

#[derive(Debug)]
pub struct Args {
    command: String,
    path: String,
    target: Option<String>,
}

impl Args {
    fn new(args: Vec<String>) -> Args {
        if args.len() != 3 && args.len() != 4 {
            panic!("invalid number of arguments in command {:?}", args)
        }
        /*
        if &args[1] != "stella" {
            panic!("stella command not used")
        }
        */
        let c = String::from(args[1].as_str());
        let p = String::from(args[2].as_str());

        let mut t = String::new();

        if args.len() == 4 {
            t = String::from(args[3].as_str());
        }

        if !t.is_empty() {
            return Args {
                command: c,
                path: p,
                target: Some(t),
            };
        }

        Args {
            command: c,
            path: p,
            target: None,
        }
    }
}

fn main() -> std::io::Result<()> {
    let args: Vec<String> = env::args().collect();
    let command_args = Args::new(args);

    let execute_ok = match command_args.command.as_str() {
        "tp" => tp::tp(&command_args),
        "new" => new::new(&command_args),
        "build" => build::build(&command_args),
        "run" => run::run(&command_args),
        _ => panic!("invalid command {}", command_args.command),
    };

    let error: Option<String> = match execute_ok {
        Ok(_) => None,
        Err(error) => Some(error.to_string()),
    };

    if error.is_some() {
        panic!(
            "error exectuing command {}: \n {}",
            command_args.command,
            error.unwrap_or(String::from("unknown error"))
        )
    }
    Ok(())
}
