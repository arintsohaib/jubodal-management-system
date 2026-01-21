"use client"

import { useState, ReactNode } from "react"
import { User, Mail, Phone, Shield, MapPin, Camera, Save, Lock, Bell, History } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { cn } from "@/pkg/utils"

export default function ProfilePage() {
    const [activeTab, setActiveTab] = useState("general")

    return (
        <div className="max-w-4xl mx-auto space-y-8">
            <div className="flex items-center gap-6">
                <div className="relative group">
                    <div className="w-24 h-24 rounded-2xl bg-gradient-to-br from-blue-600 to-emerald-600 flex items-center justify-center text-3xl font-black text-white shadow-xl group-hover:scale-105 transition-transform">
                        JD
                    </div>
                    <button className="absolute bottom-1 right-1 p-1.5 bg-background border border-border/50 rounded-lg shadow-lg opacity-0 group-hover:opacity-100 transition-opacity">
                        <Camera className="w-4 h-4 text-blue-500" />
                    </button>
                </div>
                <div className="space-y-1">
                    <h1 className="text-3xl font-black tracking-tight">Jubodal Developer</h1>
                    <p className="text-muted-foreground font-medium">Central Executive Committee • General Secretary</p>
                    <div className="flex gap-2">
                        <span className="px-2 py-0.5 bg-blue-500/10 text-blue-500 text-[10px] font-bold uppercase rounded border border-blue-500/20">Verified</span>
                        <span className="px-2 py-0.5 bg-emerald-500/10 text-emerald-500 text-[10px] font-bold uppercase rounded border border-emerald-500/20">Active</span>
                    </div>
                </div>
            </div>

            <div className="flex gap-1 p-1 bg-accent/50 rounded-xl w-fit">
                <TabButton id="general" icon={<User className="w-4 h-4" />} activeTab={activeTab} setActiveTab={setActiveTab} label="Profile" />
                <TabButton id="security" icon={<Shield className="w-4 h-4" />} activeTab={activeTab} setActiveTab={setActiveTab} label="Security" />
                <TabButton id="notifications" icon={<Bell className="w-4 h-4" />} activeTab={activeTab} setActiveTab={setActiveTab} label="Notifications" />
                <TabButton id="sessions" icon={<History className="w-4 h-4" />} activeTab={activeTab} setActiveTab={setActiveTab} label="Sessions" />
            </div>

            {activeTab === "general" && (
                <Card className="border-border/50 bg-card/50">
                    <CardHeader>
                        <CardTitle>General Information</CardTitle>
                        <CardDescription>Update your personal details and contact information.</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-6">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <Field label="Full Name (English)" value="Jubodal Developer" />
                            <Field label="নাম (বাংলা)" value="যুবদল ডেভেলপার" />
                            <Field label="Phone Number" value="+8801700000000" disabled />
                            <Field label="National ID (NID)" value="199526145678" disabled />
                        </div>
                        <div className="space-y-2">
                            <p className="text-xs font-black uppercase text-muted-foreground tracking-widest">Permanent Jurisdiction</p>
                            <div className="p-4 bg-accent/20 border border-border/50 rounded-xl flex items-center justify-between">
                                <div className="flex items-center gap-3">
                                    <MapPin className="w-5 h-5 text-blue-500" />
                                    <div>
                                        <p className="font-bold">Dhaka Division • Dhaka City South</p>
                                        <p className="text-xs text-muted-foreground">Ward 19, Unit A</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div className="flex justify-end pt-4 border-t border-border/50">
                            <Button className="bg-blue-600 hover:bg-blue-700 gap-2 font-bold">
                                <Save className="w-4 h-4" />
                                Save Changes
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}

            {activeTab === "security" && (
                <Card className="border-border/50 bg-card/50">
                    <CardHeader>
                        <CardTitle>Security Settings</CardTitle>
                        <CardDescription>Manage your password and authentication methods.</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-6">
                        <div className="flex items-center justify-between p-4 bg-accent/20 border border-border/50 rounded-xl">
                            <div className="flex items-center gap-3">
                                <Shield className="w-6 h-6 text-emerald-500" />
                                <div>
                                    <p className="font-bold">Two-Factor Authentication (2FA)</p>
                                    <p className="text-sm text-muted-foreground">Enabled via SMS (+88017***00)</p>
                                </div>
                            </div>
                            <Button variant="outline" size="sm">Manage</Button>
                        </div>
                        <div className="space-y-4">
                            <Field label="Current Password" value="••••••••" type="password" />
                            <Field label="New Password" value="" type="password" placeholder="Enter new password" />
                        </div>
                        <div className="flex justify-end pt-4 border-t border-border/50">
                            <Button className="bg-blue-600 hover:bg-blue-700 gap-2 font-bold">
                                <Lock className="w-4 h-4" />
                                Update Password
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    )
}

interface FieldProps {
    label: string
    value: string
    disabled?: boolean
    type?: string
    placeholder?: string
}

function Field({ label, value, disabled, type = "text", placeholder }: FieldProps) {
    return (
        <div className="space-y-2">
            <label className="text-xs font-black uppercase text-muted-foreground tracking-widest">{label}</label>
            <input
                type={type}
                defaultValue={value}
                disabled={disabled}
                placeholder={placeholder}
                className={cn(
                    "w-full h-11 px-4 bg-accent/20 rounded-xl border border-border/50 outline-none transition-all focus:border-blue-500/50",
                    disabled && "opacity-50 cursor-not-allowed bg-accent/10"
                )}
            />
        </div>
    )
}

interface TabButtonProps {
    id: string
    icon: ReactNode
    activeTab: string
    setActiveTab: (id: string) => void
    label: string
}

function TabButton({ id, icon, activeTab, setActiveTab, label }: TabButtonProps) {
    const active = activeTab === id
    return (
        <button
            onClick={() => setActiveTab(id)}
            className={cn(
                "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-bold transition-all",
                active ? "bg-background text-blue-500 shadow-sm" : "text-muted-foreground hover:text-foreground"
            )}
        >
            {icon}
            {label}
        </button>
    )
}
