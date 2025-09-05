async function fetchLatestPrice() {
  try {
    const res = await fetch("/crypto/latest");
    if (!res.ok) throw new Error("No data available");
    const data = await res.json();

    document.getElementById("root").innerHTML = `
      <div class="price-card">
        <p><strong>Symbol:</strong> ${data.symbol}</p>
        <p><strong>Price:</strong> $${data.price.toFixed(2)}</p>
        <p><strong>Time:</strong> ${new Date(data.time).toLocaleString()}</p>
      </div>
    `;
  } catch (err) {
    document.getElementById("root").innerHTML = `<p style="color:red;">Error: ${err.message}</p>`;
  }
}

// Initial fetch
fetchLatestPrice();
// Refresh every 5 seconds
setInterval(fetchLatestPrice, 5);
