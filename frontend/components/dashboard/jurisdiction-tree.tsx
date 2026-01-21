"use client"

import { useEffect, useState } from "react"
import { ChevronRight, ChevronDown, MapPin, Building2, Landmark, Users, Loader2 } from "lucide-react"
import { cn } from "@/pkg/utils"
import { api } from "@/pkg/api"

interface JurisdictionNode {
    id: string
    name: string
    level: string
    children?: JurisdictionNode[]
}

export function JurisdictionTree({ onSelect, activeId }: { onSelect?: (node: JurisdictionNode) => void, activeId?: string }) {
    const [hierarchy, setHierarchy] = useState<JurisdictionNode[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        async function loadHierarchy() {
            try {
                // Fetch root jurisdictions (top level)
                const res = await api.getJurisdictions()
                setHierarchy(res?.data || [])
            } catch (err) {
                console.error("Failed to load jurisdiction hierarchy:", err)
            } finally {
                setLoading(false)
            }
        }
        loadHierarchy()
    }, [])

    if (loading) {
        return (
            <div className="flex flex-col items-center justify-center py-12 opacity-50">
                <Loader2 className="w-6 h-6 animate-spin text-blue-500 mb-2" />
                <p className="text-[10px] font-bold uppercase tracking-widest">Compiling Hierarchy...</p>
            </div>
        )
    }

    return (
        <div className="space-y-2">
            <p className="text-xs font-bold uppercase tracking-widest text-muted-foreground px-4 mb-4">Organizational Hierarchy</p>
            <div className="px-2">
                {hierarchy.length > 0 ? hierarchy.map((node) => (
                    <TreeNode node={node} depth={0} onSelect={onSelect} activeId={activeId} />
                )) : (
                    <p className="text-center text-[10px] text-muted-foreground uppercase py-4">No Jurisdictions Configured</p>
                )}
            </div>
        </div>
    )
}

function TreeNode({ node, depth, onSelect, activeId }: { node: JurisdictionNode, depth: number, onSelect?: (node: JurisdictionNode) => void, activeId?: string }) {
    const [isOpen, setIsOpen] = useState(depth === 0)
    const [children, setChildren] = useState<JurisdictionNode[]>(node.children || [])
    const [loading, setLoading] = useState(false)
    const [hasLoaded, setHasLoaded] = useState(!!node.children)

    const toggleOpen = async (e: React.MouseEvent) => {
        e.stopPropagation()
        const nextState = !isOpen
        setIsOpen(nextState)

        if (nextState && !hasLoaded) {
            setLoading(true)
            try {
                const res = await api.getJurisdictions(`?parent_id=${node.id}`)
                setChildren(res?.data || [])
                setHasLoaded(true)
            } catch (err) {
                console.error("Failed to load children:", err)
            } finally {
                setLoading(false)
            }
        }
    }

    const handleSelect = () => {
        if (onSelect) onSelect(node)
    }

    const getIcon = (level: string) => {
        switch (level.toLowerCase()) {
            case 'central': return <Landmark className="w-4 h-4 text-blue-500" />
            case 'division': return <Building2 className="w-4 h-4 text-emerald-500" />
            case 'district': return <MapPin className="w-4 h-4 text-amber-500" />
            default: return <Users className="w-4 h-4 text-muted-foreground" />
        }
    }

    const leaf = !loading && hasLoaded && children.length === 0
    const isActive = activeId === node.id

    return (
        <div className="select-none">
            <div
                onClick={handleSelect}
                className={cn(
                    "flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium cursor-pointer transition-all group",
                    isActive ? "bg-blue-600 text-white shadow-lg shadow-blue-900/40" : (isOpen ? "bg-accent/50 text-foreground" : "text-muted-foreground hover:text-foreground hover:bg-accent/30"),
                    depth === 0 && !isActive && "font-bold"
                )}
            >
                <div
                    className="flex items-center justify-center w-4 h-4 hover:bg-white/10 rounded"
                    onClick={toggleOpen}
                >
                    {loading ? (
                        <Loader2 className="w-3 h-3 animate-spin" />
                    ) : (
                        !leaf && (isOpen ? <ChevronDown className="w-3 h-3 group-hover:scale-125 transition-transform" /> : <ChevronRight className="w-3 h-3 group-hover:scale-125 transition-transform" />)
                    )}
                </div>
                {getIcon(node.level)}
                <span className="truncate">{node.name}</span>
            </div>

            {isOpen && children.length > 0 && (
                <div className="ml-4 pl-4 border-l border-border/50 space-y-1 mt-1 transition-all">
                    {children.map((child) => (
                        <TreeNode key={child.id} node={child} depth={depth + 1} onSelect={onSelect} activeId={activeId} />
                    ))}
                </div>
            )}
        </div>
    )
}
