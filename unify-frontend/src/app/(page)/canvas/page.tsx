"use client";

import { useState, useCallback, ReactNode } from "react";
import {
  ReactFlow,
  applyNodeChanges,
  applyEdgeChanges,
  addEdge,
} from "@xyflow/react";

import type {
  Node,
  Edge,
  NodeChange,
  EdgeChange,
  Connection,
} from "@xyflow/react";

import "@xyflow/react/dist/style.css";
import { Server, Wifi } from "lucide-react";

import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "@/components/ui/context-menu"

// ---- Custom Node Data

type NodeData = {
  label: ReactNode;
};

// ---- Initial Nodes
const initialNodes: Node<NodeData>[] = [
  {
    id: "n1",
    position: { x: 0, y: 0 },
    data: {
      label: (

        <ContextMenu>
          <ContextMenuTrigger>
            <div style={{ display: "flex", alignItems: "center", justifyItems: "center", gap: 8 }}>
              <Wifi size={16} />
              WIFI-01
            </div>
          </ContextMenuTrigger>
          <ContextMenuContent>
            <ContextMenuItem>Profile</ContextMenuItem>
            <ContextMenuItem>Billing</ContextMenuItem>
            <ContextMenuItem>Team</ContextMenuItem>
            <ContextMenuItem>Subscription</ContextMenuItem>
          </ContextMenuContent>
        </ContextMenu>
      ),
    },
  },
  {
    id: "n2",
    position: { x: 0, y: 100 },
    data: {
      label: (

        <ContextMenu>
          <ContextMenuTrigger>
            <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
              <Server size={16} />
              SERVER-1
            </div>
          </ContextMenuTrigger>
          <ContextMenuContent>
            <ContextMenuItem>Profile</ContextMenuItem>
            <ContextMenuItem>Billing</ContextMenuItem>
            <ContextMenuItem>Team</ContextMenuItem>
            <ContextMenuItem>Subscription</ContextMenuItem>
          </ContextMenuContent>
        </ContextMenu>
      ),
    },
  },
];

const initialEdges: Edge[] = [
  {
    id: "n1-n2",
    source: "n1",
    target: "n2",
  },
];

export default function CanvassDevice() {
  const [nodes, setNodes] = useState<Node<NodeData>[]>(initialNodes);
  const [edges, setEdges] = useState<Edge[]>(initialEdges);

  const onNodesChange = useCallback((changes: NodeChange[]) => {
    setNodes((prevNodes) =>
      applyNodeChanges(changes, prevNodes) as Node<NodeData>[]
    );
  }, []);

  const onEdgesChange = useCallback((changes: EdgeChange[]) => {
    setEdges((prevEdges) =>
      applyEdgeChanges(changes, prevEdges)
    );
  }, []);

  const onConnect = useCallback((params: Connection) => {
    setEdges((prevEdges) => addEdge(params, prevEdges));
  }, []);

  return (
    <div style={{ width: "100vw", height: "100vh" }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        fitView
      />
    </div>
  );
}