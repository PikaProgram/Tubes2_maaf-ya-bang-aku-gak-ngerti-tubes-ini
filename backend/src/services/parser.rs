pub fn parse_html(html: &str) -> Result<tl::VDom<'_>, Box<dyn std::error::Error>> {
    let document = tl::parse(html, tl::ParserOptions::default())?;
    return Ok(document);
}