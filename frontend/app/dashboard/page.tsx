"use client"
import { DashboardOverview } from "@/components/dashboard/overview-widgets"

export default function DashboardPage() {
    return (
        <div className="space-y-8">
            <div>
                <h1 className="text-4xl font-black tracking-tighter">Command Center</h1>
                <p className="text-muted-foreground">Strategic overview of your jurisdiction and operations.</p>
            </div>

            <DashboardOverview />
        </div>
    )
}
