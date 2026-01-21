"use client"

import { useState } from "react"
import { Settings, Shield, Globe, Bell, Database, Save, Server, Lock } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"

export default function SettingsPage() {
    const [saving, setSaving] = useState(false)

    const handleSave = () => {
        setSaving(true)
        setTimeout(() => {
            setSaving(false)
            toast.success("Settings updated successfully")
        }, 1000)
    }

    return (
        <div className="space-y-8 pb-12">
            <div>
                <h1 className="text-4xl font-black tracking-tighter">System Control</h1>
                <p className="text-muted-foreground">Manage platform parameters, security, and jurisdictional hierarchy.</p>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
                {/* Navigation */}
                <Card className="lg:col-span-1 border-border/50 bg-card/30 h-fit sticky top-24">
                    <CardContent className="p-2 flex flex-col gap-1">
                        <button className="flex items-center gap-3 px-4 py-2 rounded-lg bg-blue-600 text-white text-sm font-bold transition-all">
                            <Globe className="w-4 h-4" /> General
                        </button>
                        <button className="flex items-center gap-3 px-4 py-2 rounded-lg text-muted-foreground hover:bg-accent/50 text-sm font-bold transition-all">
                            <Shield className="w-4 h-4" /> Permissions
                        </button>
                        <button className="flex items-center gap-3 px-4 py-2 rounded-lg text-muted-foreground hover:bg-accent/50 text-sm font-bold transition-all">
                            <Bell className="w-4 h-4" /> Alerts
                        </button>
                        <button className="flex items-center gap-3 px-4 py-2 rounded-lg text-muted-foreground hover:bg-accent/50 text-sm font-bold transition-all">
                            <Database className="w-4 h-4" /> Data Sync
                        </button>
                    </CardContent>
                </Card>

                {/* Main Content */}
                <div className="lg:col-span-3 space-y-6">
                    <Card className="border-border/50 bg-card/30 backdrop-blur-md">
                        <CardHeader>
                            <CardTitle>Organization Meta</CardTitle>
                            <CardDescription>Broad configuration for the Bangladesh Jatiotabadi Jubodal platform.</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                <div className="space-y-2">
                                    <label className="text-xs font-black uppercase tracking-widest text-muted-foreground">Organization Name</label>
                                    <input type="text" defaultValue="Bangladesh Jatiotabadi Jubodal" className="w-full h-10 px-4 rounded-lg border border-border bg-background focus:ring-1 focus:ring-blue-500 outline-none text-sm font-medium" />
                                </div>
                                <div className="space-y-2">
                                    <label className="text-xs font-black uppercase tracking-widest text-muted-foreground">System Environment</label>
                                    <div className="flex items-center gap-2 h-10 px-4 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-500 text-sm font-bold">
                                        <Server className="w-4 h-4" /> Production (grayhawks.com)
                                    </div>
                                </div>
                            </div>

                            <div className="space-y-2">
                                <label className="text-xs font-black uppercase tracking-widest text-muted-foreground">Primary Focal Point (Support Email)</label>
                                <input type="email" defaultValue="it@jubodal.org" className="w-full h-10 px-4 rounded-lg border border-border bg-background focus:ring-1 focus:ring-blue-500 outline-none text-sm font-medium" />
                            </div>
                        </CardContent>
                    </Card>

                    <Card className="border-border/50 bg-card/30 backdrop-blur-md">
                        <CardHeader>
                            <CardTitle>Security & Access Control</CardTitle>
                            <CardDescription>Configure how different ranks interact with the platform.</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <div className="flex items-center justify-between p-4 rounded-xl border border-border/50 bg-accent/20">
                                <div className="space-y-0.5">
                                    <p className="text-sm font-bold">Automatic Join Requests Audit</p>
                                    <p className="text-xs text-muted-foreground">System will automatically flag applicants with high-risk NID records.</p>
                                </div>
                                <div className="w-12 h-6 bg-blue-600 rounded-full relative p-1 cursor-pointer">
                                    <div className="w-4 h-4 bg-white rounded-full ml-auto" />
                                </div>
                            </div>

                            <div className="flex items-center justify-between p-4 rounded-xl border border-border/50 bg-accent/20">
                                <div className="space-y-0.5">
                                    <p className="text-sm font-bold">Require 2FA for District Leaders</p>
                                    <p className="text-xs text-muted-foreground">Enforces OATH/TOTP for all users at District level or higher.</p>
                                </div>
                                <div className="w-12 h-6 bg-accent rounded-full relative p-1 cursor-pointer border border-border/50">
                                    <div className="w-4 h-4 bg-muted-foreground rounded-full" />
                                </div>
                            </div>
                        </CardContent>
                    </Card>

                    <div className="flex justify-end gap-3">
                        <Button variant="outline" className="font-bold">Reset Defaults</Button>
                        <Button
                            onClick={handleSave}
                            disabled={saving}
                            className="bg-blue-600 hover:bg-blue-700 font-bold gap-2 min-w-[140px]"
                        >
                            {saving ? <Lock className="w-4 h-4 animate-spin" /> : <Save className="w-4 h-4" />}
                            Save Configuration
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    )
}
