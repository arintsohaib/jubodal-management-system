from fastapi import FastAPI, Depends, HTTPException
from fastapi.middleware.cors import CORSMiddleware
import os
import uvicorn
from pydantic import BaseModel
from typing import List, Dict, Any

from internal.intelligence import IntelligenceEngine

app = FastAPI(
    title="BJDMS Analytics AI Service",
    description="Advanced organizational intelligence for Jubodal",
    version="1.0.0"
)

# Initialize Engine
engine = IntelligenceEngine()

# CORS Configuration
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"], # In production, restrict to internal networks or specific Go-gateway
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

class HealthStatus(BaseModel):
    status: str
    service: str

@app.get("/health", response_model=HealthStatus)
def health_check():
    return {"status": "healthy", "service": "analytics"}

@app.get("/api/v1/analytics/pulse")
async def get_organizational_pulse():
    """Returns real organizational intelligence metrics."""
    try:
        growth = engine.get_growth_velocity()
        heatmap = engine.get_heatmap_data()
        performance = engine.get_unit_performance()
        
        return {
            "status": "success",
            "data": {
                "growth_velocity": growth,
                "heatmap": heatmap,
                "top_performing_units": performance
            }
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/v1/analytics/query")
async def query_assistant(payload: Dict[str, str]):
    """
    Read-only AI Assistant foundation. 
    Parses natural language queries into analytical reports.
    """
    query = payload.get("query", "").lower()
    
    if "growth" in query:
        val = engine.get_growth_velocity()
        return {"answer": f"The current organizational growth velocity is {val}% week-over-week."}
    
    if "top units" in query or "best performing" in query:
        units = engine.get_unit_performance()
        names = ", ".join([u['name'] for u in units[:3]])
        return {"answer": f"The top performing units currently are: {names}."}
        
    if "activity" in query or "heatmap" in query:
        heatmap = engine.get_heatmap_data()
        total = sum([h['activity_count'] for h in heatmap])
        return {"answer": f"System detects {total} active operations across all monitored districts."}
        
    return {"answer": "I'm sorry, I couldn't find analytical data for that specific query yet. Try asking about 'growth', 'top units', or 'activity levels'."}

if __name__ == "__main__":
    port = int(os.getenv("PORT", 8000))
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)
