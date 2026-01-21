"use client"

import { useState } from "react"
import { UserPlus, ChevronRight, ChevronLeft, MapPin, Search, FileCheck, Landmark, CheckCircle2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { cn } from "@/pkg/utils"

const steps = [
    { title: "Personal Info", icon: UserPlus },
    { title: "Location", icon: MapPin },
    { title: "Review", icon: FileCheck },
]

export default function JoinRequestWizard() {
    const [currentStep, setCurrentStep] = useState(1)
    const [submitted, setSubmitted] = useState(false)
    const [refNo, setRefNo] = useState("")

    const nextStep = () => setCurrentStep(prev => Math.min(prev + 1, steps.length))
    const prevStep = () => setCurrentStep(prev => Math.max(prev - 1, 1))

    const handleFinalSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        setRefNo("JR-2026-" + Math.floor(100000 + Math.random() * 900000))
        setSubmitted(true)
    }

    if (submitted) {
        return (
            <div className="min-h-screen flex items-center justify-center p-4 bg-background">
                <Card className="max-w-md w-full border-blue-500/20 bg-blue-500/5 backdrop-blur-md">
                    <CardContent className="pt-12 text-center space-y-6">
                        <div className="flex justify-center">
                            <div className="p-4 bg-blue-500/20 rounded-full animate-bounce">
                                <CheckCircle2 className="w-12 h-12 text-blue-500" />
                            </div>
                        </div>
                        <div className="space-y-2">
                            <h2 className="text-2xl font-bold">Application Received</h2>
                            <p className="text-muted-foreground">
                                Your membership request (Form 'Ka') has been successfully submitted to the local committee.
                            </p>
                        </div>
                        <div className="p-4 bg-accent/50 rounded-xl border border-border/50">
                            <p className="text-sm text-muted-foreground mb-1 uppercase tracking-wider font-semibold">Reference Number</p>
                            <p className="text-3xl font-mono font-bold tracking-widest text-blue-400">{refNo}</p>
                        </div>
                        <div className="text-xs text-muted-foreground space-y-1">
                            <p>• Save this number to track your status.</p>
                            <p>• You will receive an SMS once the committee reviews your application.</p>
                        </div>
                        <Button onClick={() => window.location.href = "/"} className="w-full bg-blue-600 hover:bg-blue-700">
                            Return to Home
                        </Button>
                    </CardContent>
                </Card>
            </div>
        )
    }

    return (
        <div className="min-h-screen flex flex-col bg-[radial-gradient(circle_at_bottom_right,_var(--tw-gradient-stops))] from-blue-900/10 via-background to-background">
            <header className="p-6 flex justify-between items-center border-b border-border/50 backdrop-blur-md sticky top-0 z-50">
                <div className="flex items-center gap-3">
                    <Landmark className="w-6 h-6 text-emerald-500" />
                    <span className="font-bold text-xl tracking-tight">Membership Portal</span>
                </div>
                <Button variant="ghost" onClick={() => window.location.href = "/"}>Exit</Button>
            </header>

            <main className="flex-1 flex flex-col items-center justify-center p-4 py-12">
                <div className="max-w-3xl w-full mb-12">
                    <div className="flex justify-between relative">
                        <div className="absolute top-1/2 left-0 w-full h-0.5 bg-border/50 -translate-y-1/2 -z-10" />
                        {steps.map((s, i) => (
                            <div key={i} className="flex flex-col items-center gap-2">
                                <div className={cn(
                                    "w-10 h-10 rounded-full flex items-center justify-center border-2 transition-all duration-500",
                                    currentStep > i + 1 ? "bg-emerald-500 border-emerald-500 text-white" :
                                        currentStep === i + 1 ? "bg-blue-600 border-blue-600 text-white shadow-lg shadow-blue-500/20" :
                                            "bg-background border-border text-muted-foreground"
                                )}>
                                    {currentStep > i + 1 ? <CheckCircle2 className="w-5 h-5" /> : <s.icon className="w-5 h-5" />}
                                </div>
                                <span className={cn(
                                    "text-xs font-bold uppercase tracking-wider",
                                    currentStep === i + 1 ? "text-foreground" : "text-muted-foreground"
                                )}>{s.title}</span>
                            </div>
                        ))}
                    </div>
                </div>

                <Card className="max-w-3xl w-full border-border/50 bg-card/50 backdrop-blur-xl shadow-2xl">
                    <CardHeader>
                        <CardTitle>{steps[currentStep - 1].title}</CardTitle>
                        <CardDescription>Step {currentStep} of {steps.length} — Join the movement for democracy and reform.</CardDescription>
                    </CardHeader>
                    <CardContent>
                        {currentStep === 1 && (
                            <div className="space-y-6">
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                    <div className="space-y-2">
                                        <label className="text-sm font-medium">Full Name (Bangla)</label>
                                        <input className="flex h-11 w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500" placeholder="আব্দুল করিম" />
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-sm font-medium">Full Name (English)</label>
                                        <input className="flex h-11 w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500" placeholder="Abdul Karim" />
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-sm font-medium">Phone Number</label>
                                        <input className="flex h-11 w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500" placeholder="+88017XXXXXXXX" />
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-sm font-medium">National ID (NID)</label>
                                        <input className="flex h-11 w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500" placeholder="13 digits" />
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-sm font-medium">Date of Birth</label>
                                        <input type="date" className="flex h-11 w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500" />
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-sm font-medium">Blood Group</label>
                                        <select className="flex h-11 w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500">
                                            <option>O+</option><option>O-</option><option>A+</option><option>B+</option>
                                        </select>
                                    </div>
                                </div>
                            </div>
                        )}

                        {currentStep === 2 && (
                            <div className="space-y-6">
                                <div className="space-y-2">
                                    <label className="text-sm font-medium">Target Jurisdiction (Where you want to join)</label>
                                    <select className="flex h-11 w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500">
                                        <option>Dhaka Division / Gulshan / Ward 19</option>
                                        <option>Rajshahi / Joypurhat / Ward 5</option>
                                    </select>
                                </div>
                                <div className="space-y-2">
                                    <label className="text-sm font-medium">Full Residential Address</label>
                                    <textarea className="flex min-h-[100px] w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500" placeholder="House, Road, Area..." />
                                </div>
                                <div className="space-y-2">
                                    <label className="text-sm font-medium">Referrer (Member ID or Phone - Optional)</label>
                                    <div className="relative">
                                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                                        <input className="flex h-11 w-full pl-10 rounded-md border border-input bg-background/50 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500" placeholder="Search by Member ID" />
                                    </div>
                                </div>
                            </div>
                        )}

                        {currentStep === 3 && (
                            <div className="space-y-6">
                                <div className="p-6 rounded-xl bg-blue-500/10 border border-blue-500/20 text-center space-y-2">
                                    <p className="font-bold text-lg text-blue-400">Declaration</p>
                                    <p className="text-sm text-balance">
                                        I hereby declare that I believe in the ideology of Bangladesh Nationalist Party and Jubo Dal.
                                        I am not a member of any other political organization.
                                    </p>
                                </div>
                                <div className="flex items-center space-x-2">
                                    <input type="checkbox" id="terms" className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
                                    <label htmlFor="terms" className="text-sm text-muted-foreground">
                                        I agree to the <span className="text-blue-500 hover:underline cursor-pointer">Constitutional Rules</span> of the organization.
                                    </label>
                                </div>
                            </div>
                        )}
                    </CardContent>
                    <CardFooter className="flex justify-between border-t border-border/50 pt-6">
                        <Button
                            variant="ghost"
                            onClick={prevStep}
                            disabled={currentStep === 1}
                            className="h-11 px-6"
                        >
                            <ChevronLeft className="w-4 h-4 mr-2" />
                            Previous
                        </Button>
                        {currentStep < steps.length ? (
                            <Button onClick={nextStep} className="bg-blue-600 hover:bg-blue-700 h-11 px-8 font-semibold">
                                Next Step
                                <ChevronRight className="w-4 h-4 ml-2" />
                            </Button>
                        ) : (
                            <Button onClick={handleFinalSubmit} className="bg-emerald-600 hover:bg-emerald-700 h-11 px-10 font-bold shadow-lg shadow-emerald-500/20">
                                Submit Application
                            </Button>
                        )}
                    </CardFooter>
                </Card>
            </main>

            <footer className="p-8 border-t border-border/50 text-center text-sm text-muted-foreground bg-accent/20">
                <p>Membership Fee (Tk 10) is collected after committee approval during verification.</p>
            </footer>
        </div>
    )
}
