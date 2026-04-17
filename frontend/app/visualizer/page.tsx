"use client";

import { motion } from "motion/react";
import Link from "next/link";
import React from "react";
import { TraversalApiResponse, TraversalRequestPayload } from "@/types/traversal";

function parseSessionJSON<T>(value: string | null): T | null {
  if (!value) {
    return null;
  }

  try {
    return JSON.parse(value) as T;
  } catch {
    return null;
  }
}

export default function VisualizerPage() {
  const [requestPayload, setRequestPayload] = React.useState<TraversalRequestPayload | null>(null);
  const [responsePayload, setResponsePayload] = React.useState<TraversalApiResponse | null>(null);
  const [isSessionChecked, setIsSessionChecked] = React.useState(false);

  React.useEffect(() => {
    const savedRequest = parseSessionJSON<TraversalRequestPayload>(
      sessionStorage.getItem("traversal:request")
    );
    const savedResponse = parseSessionJSON<TraversalApiResponse>(
      sessionStorage.getItem("traversal:response")
    );

    setRequestPayload(savedRequest);
    setResponsePayload(savedResponse);
    setIsSessionChecked(true);
  }, []);

  if (!isSessionChecked) {
    return (
      <div className="min-h-screen bg-background text-white selection:bg-primary selection:text-on-primary px-6 py-12">
        <div className="max-w-3xl mx-auto border border-outline-variant/40 bg-surface-container/40 p-8">
          <h1 className="text-4xl font-black text-white mb-8 leading-[0.9] tracking-tighter uppercase max-w-full">
            <span className="text-primary italic drop-shadow-[0_0_32px_rgba(100,240,255,0.5)]">
              Traversal Dashboard
            </span>
          </h1>
          <p className="mt-4 text-on-surface-variant">
            Loading traversal data...
          </p>
        </div>
      </div>
    );
  }

  if (!responsePayload) {
    return (
      <div className="min-h-screen bg-background text-white selection:bg-primary selection:text-on-primary px-6 py-12">
        <div className="max-w-3xl mx-auto border border-outline-variant/40 bg-surface-container/40 p-8">
          <h1 className="text-4xl font-black text-white mb-8 leading-[0.9] tracking-tighter uppercase max-w-full">
            <span className="text-primary italic drop-shadow-[0_0_32px_rgba(100,240,255,0.5)]">
              Traversal Dashboard
            </span>
          </h1>
          <p className="mt-4 text-on-surface-variant">
            No traversal data found. Submit the form from homepage first.
          </p>
          <Link
            href="/"
            className="inline-flex mt-6 bg-primary text-on-primary font-black px-5 py-3 uppercase tracking-wider"
          >
            Back To Homepage
          </Link>
        </div>
      </div>
    );
  }

  return (<div>Data Available</div>);
}