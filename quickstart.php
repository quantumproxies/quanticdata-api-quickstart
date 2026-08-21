<?php
/**
 * QuanticData API quickstart — PHP (curl extension, no Composer).
 *   QUANTICDATA_API_KEY=qd_live_... php quickstart.php
 */
declare(strict_types=1);

const QD_BASE = 'https://api.quanticdata.io/v1';

/** POST to the API and unwrap the { type, message, payload } envelope. */
function qd_call(string $path, array $body): array
{
    $key = getenv('QUANTICDATA_API_KEY');
    $ch = curl_init(QD_BASE . $path);
    curl_setopt_array($ch, [
        CURLOPT_POST           => true,
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_TIMEOUT        => 120,
        CURLOPT_HTTPHEADER     => [
            'Authorization: Bearer ' . $key,
            'Content-Type: application/json',
        ],
        CURLOPT_POSTFIELDS     => json_encode($body, JSON_UNESCAPED_SLASHES),
    ]);
    $raw    = curl_exec($ch);
    $status = curl_getinfo($ch, CURLINFO_RESPONSE_CODE);
    curl_close($ch);

    $data = json_decode((string) $raw, true) ?: [];
    if ($status >= 400 || ($data['type'] ?? '') === 'error') {
        throw new RuntimeException(sprintf('%s failed (%d): %s', $path, $status, $data['message'] ?? 'unknown'));
    }
    return $data['payload'] ?? [];
}

if (!getenv('QUANTICDATA_API_KEY')) {
    fwrite(STDERR, "set QUANTICDATA_API_KEY — https://app.quanticdata.io/register\n");
    exit(1);
}

$page = qd_call('/scrape', ['url' => 'https://example.com', 'format' => 'markdown']);
echo "--- markdown -------------------------------------------------\n";
echo substr($page['markdown'] ?? '', 0, 400), "\n";
echo 'title: ', $page['metadata']['title'] ?? '', "\n";
echo 'cost : ', $page['usage']['cost_usd'] ?? '', "\n";

$serp = qd_call('/serp', ['query' => 'web scraping api', 'country' => 'us', 'num' => 5]);
echo "\n--- organic --------------------------------------------------\n";
foreach (array_slice($serp['organic'] ?? [], 0, 5) as $row) {
    printf("%2d. %s\n    %s\n", $row['position'] ?? 0, $row['title'] ?? '', $row['link'] ?? '');
}
