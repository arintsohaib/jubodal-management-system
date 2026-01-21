import { Button } from "@/components/ui/button";
import Link from "next/link";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Shield, Users, Activity, Landmark, Flag } from "lucide-react";

export default function Home() {
    return (
        <div className="flex flex-col items-center justify-center py-20 px-4">
            <div className="max-w-4xl w-full space-y-12 text-center">
                <div className="space-y-4">
                    <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight bg-gradient-to-r from-blue-400 to-emerald-400 bg-clip-text text-transparent">
                        BJDMS
                    </h1>
                    <p className="text-xl md:text-2xl text-muted-foreground max-w-2xl mx-auto">
                        The next-generation management platform for Bangladesh Jatiotabadi Jubodal.
                        Auditable, Transparent, and Secure.
                    </p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    <FeatureCard
                        icon={<Shield className="w-8 h-8 text-blue-400" />}
                        title="Secure Auth"
                        description="JWT-based security with hierarchical access control."
                    />
                    <FeatureCard
                        icon={<Users className="w-8 h-8 text-emerald-400" />}
                        title="Committees"
                        description="Grassroots to Central committee management."
                    />
                    <FeatureCard
                        icon={<Activity className="w-8 h-8 text-orange-400" />}
                        title="Activity"
                        description="Daily logging and task tracking for all members."
                    />
                    <FeatureCard
                        icon={<Landmark className="w-8 h-8 text-purple-400" />}
                        title="Finance"
                        description="Immutable donation and expense ledger."
                    />
                </div>

                <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
                    <Link href="/public/join">
                        <button className="px-8 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold transition-all shadow-lg shadow-blue-900/40">
                            Apply to Join
                        </button>
                    </Link>
                    <Link href="/dashboard">
                        <button className="px-8 py-3 bg-muted hover:bg-muted/80 text-foreground rounded-lg font-semibold transition-all">
                            Login to Dashboard
                        </button>
                    </Link>
                </div>
            </div>
        </div>
    );
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode, title: string, description: string }) {
    return (
        <div className="p-6 rounded-2xl bg-accent/50 border border-border/50 backdrop-blur-sm hover:border-blue-500/50 transition-colors text-left space-y-4">
            <div className="p-3 bg-background rounded-xl w-fit">
                {icon}
            </div>
            <h3 className="text-xl font-bold">{title}</h3>
            <p className="text-sm text-muted-foreground leading-relaxed">
                {description}
            </p>
        </div>
    );
}
