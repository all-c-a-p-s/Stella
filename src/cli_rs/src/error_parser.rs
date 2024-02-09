pub fn parse_error(error_msg: String) -> String {
    let parts = error_msg.as_str().split('\n');
    let lines = parts.collect::<Vec<&str>>();
    if lines.is_empty() {
        panic!("error message is empty")
    }
    String::from(lines[0])
}
