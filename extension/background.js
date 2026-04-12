const API = "https://postificus-api.onrender.com";

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

    const hasCreds = platform === "devto"
        ? creds["remember_user_token"]
        : creds["uid"] && creds["sid"];

    if (!hasCreds) return;

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

function checkUrl(url) {
    if (!url) return;
    if (url.includes("dev.to")) syncPlatform("devto", ["remember_user_token"], "dev.to");
    else if (url.includes("medium.com")) syncPlatform("medium", ["uid", "sid", "xsrf"], "medium.com");
}

// Trigger on tab navigation
chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
    if (changeInfo.status === "complete") checkUrl(tab.url);
});

// Trigger immediately on install/reload for already-open tabs
chrome.runtime.onInstalled.addListener(() => {
    chrome.tabs.query({}, (tabs) => tabs.forEach(t => checkUrl(t.url)));
});

// Also trigger on browser startup
chrome.runtime.onStartup.addListener(() => {
    chrome.tabs.query({}, (tabs) => tabs.forEach(t => checkUrl(t.url)));
});
