const BASE = ''

export async function get<T>(url: string): Promise<T> {
  const res = await fetch(BASE + url)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}