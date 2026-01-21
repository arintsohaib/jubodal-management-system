"use client"

import { useEffect, useState } from "react"
import { useAuth } from "@/pkg/auth-context"
import { toast } from "sonner"

export interface Notification {
    id: string
    type: string
    title: string
    message: string
    data?: any
    is_read: boolean
    created_at: string
}

export function useNotifications() {
    const { user } = useAuth()
    const [notifications, setNotifications] = useState<Notification[]>([])
    const [unreadCount, setUnreadCount] = useState(0)

    useEffect(() => {
        if (!user) return

        // Initial Fetch
        fetch(`${process.env.NEXT_PUBLIC_API_URL}/notifications`, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        })
            .then(res => res.json())
            .then(data => {
                if (Array.isArray(data)) {
                    setNotifications(data)
                    setUnreadCount(data.filter((n: any) => !n.is_read).length)
                }
            })

        // WebSocket Connection
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const wsUrl = `${protocol}//${window.location.host}/api/v1/notifications/ws`

        const socket = new WebSocket(wsUrl)

        socket.onmessage = (event) => {
            const notification: Notification = JSON.parse(event.data)
            setNotifications(prev => [notification, ...prev])
            setUnreadCount(prev => prev + 1)

            // Show real-time toast
            toast(notification.title, {
                description: notification.message,
            })
        }

        return () => {
            socket.close()
        }
    }, [user])

    const markAsRead = async (id: string) => {
        await fetch(`${process.env.NEXT_PUBLIC_API_URL}/notifications/${id}/read`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        })

        setNotifications(prev => prev.map(n =>
            n.id === id ? { ...n, is_read: true } : n
        ))
        setUnreadCount(prev => Math.max(0, prev - 1))
    }

    return {
        notifications,
        unreadCount,
        markAsRead
    }
}
