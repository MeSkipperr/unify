'use client'

import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
} from '@/components/ui/sidebar'
import Link from 'next/link';


export default function AppSidebar() {
    const property = process.env.NEXT_PUBLIC_PROPERTY;
    const date = new Date(); 
    return (
        <Sidebar  >
            <SidebarContent className='bg-background'>
                <SidebarGroup>
                    <SidebarGroupLabel className='text-xl py-6 font-bold '>Menu</SidebarGroupLabel>
                    <SidebarGroupContent>
                        <SidebarMenu>

                            {/* Dashboard */}
                            <SidebarMenuItem>
                                <SidebarMenuButton asChild className="font-bold">
                                    <Link href="/dashboard">Dashboard</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            {/* Devices */}
                            <SidebarMenuItem>
                                <SidebarMenuButton asChild className="font-bold">
                                    <Link href="/devices">Devices</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem className="pl-4">
                                <SidebarMenuButton asChild>
                                    <Link href="/devices/access-point">Access Point</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem className="pl-4">
                                <SidebarMenuButton asChild>
                                    <Link href="/devices/cctv">CCTV</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem className="pl-4">
                                <SidebarMenuButton asChild>
                                    <Link href="/devices/iptv">IPTV</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            {/* Services */}
                            <SidebarMenuItem>
                                <SidebarMenuButton asChild className="font-bold">
                                    <Link href="/services">Services</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem className="pl-4">
                                <SidebarMenuButton asChild>
                                    <Link href="/services/network-monitoring">
                                        Monitoring Network
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem className="pl-4">
                                <SidebarMenuButton asChild>
                                    <Link href="/services/adb">
                                        Android Debug Bridge
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem className="pl-8">
                                <SidebarMenuButton asChild>
                                    <Link href="/services/adb/remove-youtube-data">
                                        Remove Youtube Data
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem className="pl-8">
                                <SidebarMenuButton asChild>
                                    <Link href="/services/adb/uptime">
                                        Get Uptime
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem className="pl-4">
                                <SidebarMenuButton asChild>
                                    <Link href="/services/speed-test">Speed Test</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem className="pl-4">
                                <SidebarMenuButton asChild>
                                    <Link href="/services/port-forward">Port Forward</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem className="pl-4">
                                <SidebarMenuButton asChild>
                                    <Link href="/services/traceroute">Trace Route</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            {/* Others */}
                            <SidebarMenuItem>
                                <SidebarMenuButton asChild className="font-bold">
                                    <Link href="/logs">Logs</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem>
                                <SidebarMenuButton asChild className="font-bold">
                                    <Link href="/users">Users</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem>
                                <SidebarMenuButton asChild className="font-bold">
                                    <Link href="/system">System</Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                        </SidebarMenu>
                    </SidebarGroupContent>


                </SidebarGroup>
            </SidebarContent>


            <SidebarFooter className="text-sm text-muted-foreground bg-background border-t">

                {property} © Unify {date.getFullYear()}
            </SidebarFooter>
        </Sidebar>
    )
}