from sqlalchemy import create_engine, text
import pandas as pd
import os
from typing import List, Dict, Any

class IntelligenceEngine:
    def __init__(self):
        self.db_url = os.getenv("DATABASE_URL")
        self.engine = create_engine(self.db_url)

    def get_growth_velocity(self) -> float:
        """
        Calculates the week-over-week growth percentage of new member signups.
        """
        query = text("""
            WITH weekly_counts AS (
                SELECT date_trunc('week', created_at) as week, count(*) as member_count
                FROM users
                GROUP BY 1
                ORDER BY 1 DESC
                LIMIT 2
            )
            SELECT member_count FROM weekly_counts
        """)
        
        with self.engine.connect() as conn:
            results = conn.execute(query).fetchall()
            
        if len(results) < 2:
            return 0.0
            
        latest = results[0][0]
        previous = results[1][0]
        
        if previous == 0:
            return 100.0
            
        return round(((latest - previous) / previous) * 100, 2)

    def get_heatmap_data(self) -> List[Dict[str, Any]]:
        """
        Aggregates activity density by district for heatmap visualization.
        """
        query = text("""
            SELECT j.name as district_name, count(a.id) as activity_count
            FROM activities a
            JOIN jurisdictions j ON a.jurisdiction_id = j.id
            WHERE j.level = 'district'
            GROUP BY j.name
            ORDER BY activity_count DESC
        """)
        
        with self.engine.connect() as conn:
            df = pd.read_sql(query, conn)
            
        return df.to_dict(orient="records")

    def get_unit_performance(self) -> List[Dict[str, Any]]:
        """
        Calculates performance scores for competitive unit leaderboards.
        """
        query = text("""
            SELECT j.name, count(a.id) * 10 + count(u.id) as score
            FROM jurisdictions j
            LEFT JOIN activities a ON a.jurisdiction_id = j.id
            LEFT JOIN users u ON u.jurisdiction_id = j.id
            WHERE j.level IN ('district', 'upazila')
            GROUP BY j.name
            ORDER BY score DESC
            LIMIT 10
        """)
        
        with self.engine.connect() as conn:
            df = pd.read_sql(query, conn)
            
        return df.to_dict(orient="records")
