pub async fn fetch_html_page(url: &str) -> Result<String, Box<dyn std::error::Error>> {
    let normalized = if url.contains("://") {
        url.to_string()
    } else {
        format!("https://{}", url)
    };

    let parsed = reqwest::Url::parse(&normalized)?;
    let response = reqwest::get(parsed).await?;
    let html = response.text().await?;
    Ok(html)
}