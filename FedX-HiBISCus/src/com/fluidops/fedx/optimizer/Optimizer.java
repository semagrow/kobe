/*
 * Copyright (C) 2008-2012, fluid Operations AG
 *
 * FedX is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 * 
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package com.fluidops.fedx.optimizer;

import java.io.IOException;
import java.util.List;
import java.util.Map;

import  org.aksw.simba.fedsum.startup.QueryEvaluation;

import java.util.Set;

import org.aksw.simba.fedsum.FedSumSourceSelection;
import org.apache.log4j.Logger;
import org.openrdf.query.BindingSet;
import org.openrdf.query.Dataset;
import org.openrdf.query.MalformedQueryException;
import org.openrdf.query.QueryEvaluationException;
import org.openrdf.query.algebra.QueryRoot;
import org.openrdf.query.algebra.StatementPattern;
import org.openrdf.query.algebra.TupleExpr;
import org.openrdf.query.algebra.evaluation.impl.ConstantOptimizer;
import org.openrdf.query.algebra.evaluation.impl.DisjunctiveConstraintOptimizer;
import org.openrdf.query.algebra.evaluation.impl.EvaluationStrategyImpl;
import org.openrdf.repository.RepositoryException;
import org.openrdf.sail.SailException;

import com.fluidops.fedx.FedX;
import com.fluidops.fedx.FederationManager;
import com.fluidops.fedx.algebra.SingleSourceQuery;
import com.fluidops.fedx.algebra.StatementSource;
import com.fluidops.fedx.cache.Cache;
import com.fluidops.fedx.structures.Endpoint;
import com.fluidops.fedx.structures.QueryInfo;



public class Optimizer {

	public static Logger logger = Logger.getLogger(Optimizer.class);
	
	
	public static TupleExpr optimize(TupleExpr parsed, Dataset dataset, BindingSet bindings, 
			EvaluationStrategyImpl strategy, QueryInfo queryInfo) throws SailException {
		
		FedX fed = FederationManager.getInstance().getFederation();
		List<Endpoint> members = fed.getMembers();
		
		// if the federation has a single member only, evaluate the entire query there
		if (members.size()==1 && queryInfo.getQuery()!=null)
			return new SingleSourceQuery(parsed, members.get(0), queryInfo);			
		
		// Clone the tuple expression to allow for more aggressive optimizations
		TupleExpr query = new QueryRoot(parsed.clone());
		
		Cache cache = FederationManager.getInstance().getCache();

		if (logger.isTraceEnabled())
			logger.trace("Query before Optimization: " + query);
		
		
		/* original sesame optimizers */
		new ConstantOptimizer(strategy).optimize(query, dataset, bindings);		// maybe remove this optimizer later

		new DisjunctiveConstraintOptimizer().optimize(query, dataset, bindings);

		
		/*
		 * TODO
		 * add some generic optimizers: 
		 *  - FILTER ?s=1 && ?s=2 => EmptyResult
		 *  - Remove variables that are not occuring in query stmts from filters
		 */
		
		
		/* custom optimizers, execute only when needed*/
			
		GenericInfoOptimizer info = new GenericInfoOptimizer(queryInfo);
		
		// collect information and perform generic optimizations
		info.optimize(query);
		
		// Source Selection: all nodes are annotated with their source
		long srcTime = System.currentTimeMillis();
		SourceSelection sourceSelection = new SourceSelection(members, cache, queryInfo);
		try {
			sourceSelection.doSourceSelection(info.getStatements());
			long srcSelTime = System.currentTimeMillis()-srcTime;
		    System.out.println("Source Selection Time "+ (srcSelTime) );
		    try {
				QueryEvaluation.bw.write(""+srcSelTime);
			} catch (IOException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			}
			  int tpsrces = 0; 
				for (StatementPattern stmt : SourceSelection.stmtToSources.keySet()) 
		          {
		        	tpsrces = tpsrces+ SourceSelection.stmtToSources.get(stmt).size();
//					System.out.println("-----------\n"+stmt);
//					System.out.println(SourceSelection.stmtToSources.get(stmt));
		         }
			      System.out.println("Total Triple pattern-wise selected sources after step 2 of FedSum source selection : "+ tpsrces);
			  
	
		} catch (RepositoryException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (MalformedQueryException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (QueryEvaluationException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
//		System.out.println("Source Selection Time "+ (System.currentTimeMillis()-srcTime) );	
//		
//		FedSumSourceSelection fsourceSelection = new FedSumSourceSelection(members,cache,queryInfo);
//		try {
//			long FedSumTime = System.currentTimeMillis();
//			  Map<StatementPattern, List<StatementSource>> stmtToSources=	fsourceSelection.performSourceSelection(QueryEvaluation.DNFgrps);
//			System.out.println("Source Selection Time "+ (System.currentTimeMillis()-FedSumTime) );	
//			  int tpsrces = 0; 
//				for (StatementPattern stmt : stmtToSources.keySet()) 
//		          {
//		        	tpsrces = tpsrces+ stmtToSources.get(stmt).size();
//					System.out.println("-----------\n"+stmt);
//					System.out.println(stmtToSources.get(stmt));
//		         }
//			      System.out.println("Total Triple pattern-wise selected sources after step 2 of FedSum source selection : "+ tpsrces);
//			
//		} catch (RepositoryException e) {
//			// TODO Auto-generated catch block
//			e.printStackTrace();
//		} catch (MalformedQueryException e) {
//			// TODO Auto-generated catch block
//			e.printStackTrace();
//		} catch (QueryEvaluationException e) {
//			// TODO Auto-generated catch block
//			e.printStackTrace();
//		}
		// if the query has a single relevant source (and if it is no a SERVICE query), evaluate at this source only
		Set<Endpoint> relevantSources = sourceSelection.getRelevantSources();
		if (relevantSources.size()==1 && !info.hasService())
			return new SingleSourceQuery(query, relevantSources.iterator().next(), queryInfo);		
		
		if (info.hasService())
			new ServiceOptimizer(queryInfo).optimize(query);
		
		// optimize Filters, if available
		if (info.hasFilter())
			new FilterOptimizer().optimize(query);

		// optimize unions, if available
		if (info.hasUnion)
			new UnionOptimizer(queryInfo).optimize(query);
		
		// optimize statement groups and join order
		new StatementGroupOptimizer(queryInfo).optimize(query);
				
		new VariableScopeOptimizer(queryInfo).optimize(query);
		
		if (logger.isTraceEnabled())
			logger.trace("Query after Optimization: " + query);

		return query;
	}	

}