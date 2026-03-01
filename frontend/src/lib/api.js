const DEFAULT_BASE = import.meta.env.VITE_API_BASE || '';

export async function get(path, opts = {}) {
  const res = await fetch(`${DEFAULT_BASE}${path}`, {
    credentials: 'same-origin',
    ...opts,
  });
  return res;
}

export default { get };
