"use client"

import { useState } from "react"
import { ShieldAlert, Send, FileText, CheckCircle2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { cn } from "@/pkg/utils"

export default function ComplaintPublicPage() {
    const [step, setStep] = useState(1)
    const [isAnonymous, setIsAnonymous] = useState(true)
    const [submitted, setSubmitted] = useState(false)
    const [trackId, setTrackId] = useState("")

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        // Mock submission
        setTrackId("C-" + Math.random().toString(36).substr(2, 6).toUpperCase())
        setSubmitted(true)
    }

    if (submitted) {
        return (
            <div className="min-h-screen flex items-center justify-center p-4 bg-background">
                <Card className="max-w-md w-full border-emerald-500/20 bg-emerald-500/5 backdrop-blur-md">
                    <CardContent className="pt-12 text-center space-y-6">
                        <div className="flex justify-center">
                            <div className="p-4 bg-emerald-500/20 rounded-full">
                                <CheckCircle2 className="w-12 h-12 text-emerald-500" />
                            </div>
                        </div>
                        <div className="space-y-2">
                            <h2 className="text-2xl font-bold">Complaint Submitted</h2>
                            <p className="text-muted-foreground">
                                Thank you for your feedback. We take all reports seriously and will investigate accordingly.
                            </p>
                        </div>
                        <div className="p-4 bg-accent/50 rounded-xl border border-border/50">
                            <p className="text-sm text-muted-foreground mb-1 uppercase tracking-wider font-semibold">Your Tracking ID</p>
                            <p className="text-3xl font-mono font-bold tracking-widest text-emerald-400">{trackId}</p>
                        </div>
                        <Button onClick={() => window.location.href = "/"} className="w-full">
                            Back to Home
                        </Button>
                    </CardContent>
                </Card>
            </div>
        )
    }

    return (
        <div className="min-h-screen flex flex-col bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-red-900/10 via-background to-background">
            <header className="p-6 flex justify-between items-center border-b border-border/50 backdrop-blur-md sticky top-0 z-50">
                <div className="flex items-center gap-3">
                    <ShieldAlert className="w-6 h-6 text-red-500" />
                    <span className="font-bold text-xl tracking-tight">Public Complaint Portal</span>
                </div>
                <Button variant="ghost" onClick={() => window.location.href = "/"}>Cancel</Button>
            </header>

            <main className="flex-1 flex items-center justify-center p-4 py-12">
                <Card className="max-w-2xl w-full border-border/50 bg-card/50 backdrop-blur-xl shadow-2xl">
                    <CardHeader>
                        <CardTitle className="text-2xl">Submit a Report</CardTitle>
                        <CardDescription>
                            Provide details about the issue or grievance. You can choose to remain completely anonymous.
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleSubmit} className="space-y-8">
                            <div className="flex p-1 bg-accent/50 rounded-lg w-fit">
                                <button
                                    type="button"
                                    onClick={() => setIsAnonymous(true)}
                                    className={cn(
                                        "px-4 py-2 rounded-md text-sm font-medium transition-all",
                                        isAnonymous ? "bg-background shadow-sm text-foreground" : "text-muted-foreground"
                                    )}
                                >
                                    Anonymous
                                </button>
                                <button
                                    type="button"
                                    onClick={() => setIsAnonymous(false)}
                                    className={cn(
                                        "px-4 py-2 rounded-md text-sm font-medium transition-all",
                                        !isAnonymous ? "bg-background shadow-sm text-foreground" : "text-muted-foreground"
                                    )}
                                >
                                    Identified
                                </button>
                            </div>

                            <div className="space-y-6">
                                {!isAnonymous && (
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                        <div className="space-y-2">
                                            <label className="text-sm font-medium">Full Name</label>
                                            <input className="flex h-10 w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" placeholder="John Doe" />
                                        </div>
                                        <div className="space-y-2">
                                            <label className="text-sm font-medium">Phone Number</label>
                                            <input className="flex h-10 w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" placeholder="+88017XXXXXXXX" />
                                        </div>
                                    </div>
                                )}

                                <div className="space-y-2">
                                    <label className="text-sm font-medium">Jurisdiction / Location of Incident</label>
                                    <select className="flex h-10 w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
                                        <option>Select Jurisdiction</option>
                                        <option>Dhaka Division</option>
                                        <option>Chittagong Division</option>
                                        <option>Rajshahi Division</option>
                                    </select>
                                </div>

                                <div className="space-y-2">
                                    <label className="text-sm font-medium">Complaint Details</label>
                                    <textarea
                                        className="flex min-h-[120px] w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                                        placeholder="Provide a clear description of the incident, including date, time, and involved parties..."
                                        required
                                    />
                                </div>

                                <div className="space-y-2">
                                    <label className="text-sm font-medium">Evidence (Optional)</label>
                                    <div className="border-2 border-dashed border-border/50 rounded-xl p-8 text-center space-y-2 hover:border-blue-500/50 transition-colors cursor-pointer bg-background/20">
                                        <FileText className="w-8 h-8 mx-auto text-muted-foreground" />
                                        <p className="text-sm text-muted-foreground">Drag and drop images or documents, or click to browse</p>
                                        <p className="text-xs text-muted-foreground/50">Max size: 5MB. Supported: JPG, PNG, PDF</p>
                                    </div>
                                </div>
                            </div>

                            <div className="flex justify-end gap-4 pt-4">
                                <Button type="submit" className="bg-red-600 hover:bg-red-700 h-11 px-8">
                                    Submit Report
                                    <Send className="w-4 h-4 ml-2" />
                                </Button>
                            </div>
                        </form>
                    </CardContent>
                </Card>
            </main>

            <footer className="p-8 border-t border-border/50 text-center text-sm text-muted-foreground bg-accent/20">
                <p>&copy; 2026 Bangladesh Jatiotabadi Jubodal. Secure, Encrypted, and Auditable.</p>
            </footer>
        </div>
    )
}
