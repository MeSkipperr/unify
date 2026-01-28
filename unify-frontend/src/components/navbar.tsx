'use client'


import { Bell, LogOut, User } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { SidebarTrigger } from '@/components/ui/sidebar'
import { ModeToggle } from './mode-toggle'


export default function AppNavbar() {
    const newNotification = true;
    return (
        <nav className="flex h-14 items-center justify-between border-b px-4 sticky top-0 bg-background">
            <div className="flex items-center gap-2">
                <SidebarTrigger className='size-10 cursor-pointer' />
                <h1 className="text-xl font-extrabold tracking-wide text-primary">Unify</h1>
            </div>


            <div className="flex items-center gap-2">
                <Button
                    variant="ghost"
                    size="icon"
                    className="relative hover:cursor-pointer p-2 "
                >
                    <Bell className="size-5" />

                    {newNotification &&
                        <span className="absolute top-1 right-1 h-1.5 w-1.5 rounded-full bg-red-500" />
                    }
                </Button>

                <ModeToggle />


                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon" className="hover:cursor-pointer p-2">
                            <User className="size-6" />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        <DropdownMenuItem className="text-destructive hover:cursor-pointer">
                            <LogOut className="mr-2 h-4 w-4" /> Logout
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </nav>
    )
}