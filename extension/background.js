const API = "https://postificus-api.onrender.com";
const USER_ID = "00000000-0000-0000-0000-000000000001";

async function syncPlatform(platform, cookieNames, domain) {
    const results = await Promise.all(
        cookieNames.map(name =>
            chrome.cookies.get({ url: `https://${domain}`, name })
        )
    );

    const creds = {};
    results.forEach((cookie, i) => {
        if (cookie?.value) creds[cookieNames[i]] = cookie.value;
    });

    // Only send if we have at least the primary credential
    const hasCreds = platform === "devto"
        ? creds["remember_user_token"]
        : creds["uid"] && creds["sid"];

    if (!hasCreds) return;

    // Check if already saved (avoid hammering the API on every page load)
    const stored = await chrome.storage.local.get(`synced_${platform}`);
    const key = JSON.stringify(creds);
    if (stored[`synced_${platform}`] === key) return;

    await fetch(`${API}/api/settings/credentials`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ platform, credentials: creds })
    });

    await chrome.storage.local.set({ [`synced_${platform}`]: key });
    console.log(`[Postificus] ✅ ${platform} credentials synced`);
}

chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
    if (changeInfo.status !== "complete") return;

    if (tab.url?.includes("dev.to")) {
        syncPlatform("devto", ["remember_user_token"], "dev.to");
    } else if (tab.url?.includes("medium.com")) {
        syncPlatform("medium", ["uid", "sid", "xsrf"], "medium.com");
    }
});
