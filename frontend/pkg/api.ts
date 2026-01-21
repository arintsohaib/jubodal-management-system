const API_URL = (typeof process !== "undefined" && process.env.NEXT_PUBLIC_API_URL) || "https://arint.win/api/v1";

async function fetchWithAuth(endpoint: string, options: RequestInit = {}) {
    const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;
    const headers = {
        "Content-Type": "application/json",
        ...options.headers,
        ...(token ? { "Authorization": `Bearer ${token}` } : {}),
    } as any;

    const response = await fetch(`${API_URL}${endpoint}`, { ...options, headers });

    if (response.status === 401) {
        if (typeof window !== "undefined") window.location.href = "/login";
        return null;
    }

    if (!response.ok) {
        throw new Error(`API Error: ${response.statusText}`);
    }

    return response.json();
}

export const api = {
    // Auth
    getCurrentUser: () => fetchWithAuth("/auth/me"),

    // Analytics
    getPulse: () => fetchWithAuth("/analytics/pulse"),

    // Join Requests
    getJoinRequests: (params: string = "") => fetchWithAuth(`/join-requests${params}`),
    approveJoinRequest: (id: string) => fetchWithAuth(`/join-requests/${id}/approve`, { method: "PATCH" }),
    rejectJoinRequest: (id: string, reason: string) => fetchWithAuth(`/join-requests/${id}/reject`, {
        method: "PATCH",
        body: JSON.stringify({ reason })
    }),

    // Activities
    getActivities: (params: string = "") => fetchWithAuth(`/activities${params}`),
    logActivity: (data: any) => fetchWithAuth("/activities", {
        method: "POST",
        body: JSON.stringify(data)
    }),

    // Organizations
    getCommittees: () => fetchWithAuth("/org/committees"),
    getJurisdictions: (params: string = "") => fetchWithAuth(`/org/jurisdictions${params}`),
    getCommitteeMembers: (committeeId: string) => fetchWithAuth(`/org/committees/${committeeId}/members`),
    addCommitteeMember: (committeeId: string, data: any) => fetchWithAuth(`/org/committees/${committeeId}/members`, {
        method: "POST",
        body: JSON.stringify(data)
    }),

    // Finance
    getFinanceStatement: (jurisdictionId: string, page: number = 1) =>
        fetchWithAuth(`/finance/statement?jurisdiction_id=${jurisdictionId}&page=${page}`),
    recordTransaction: (data: any) => fetchWithAuth("/finance/transactions", {
        method: "POST",
        body: JSON.stringify(data)
    }),

    // Search
    search: (query: string) => fetchWithAuth(`/search?q=${encodeURIComponent(query)}`),
};
