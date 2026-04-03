#[path = "services/scraper.rs"]
mod scraper;
#[path = "services/parser.rs"]
mod parser;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let raw_data = scraper::fetch_html_page("x.com").await?;
    println!("{}", raw_data);

    let parsed_data = parser::parse_html(&raw_data)?;
    println!("{:#?}", parsed_data);

    Ok(())
}
