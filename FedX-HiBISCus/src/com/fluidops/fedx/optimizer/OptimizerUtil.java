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

import java.util.ArrayList;
import java.util.List;

import org.openrdf.query.algebra.Join;
import org.openrdf.query.algebra.TupleExpr;

import com.fluidops.fedx.algebra.NJoin;
import com.fluidops.fedx.structures.QueryInfo;

public class OptimizerUtil
{

	
	/**
	 * Flatten the join to one layer, i.e. collect all join arguments
	 * 
	 * @param join
	 * @param queryInfo
	 * @return
	 */
	public static NJoin flattenJoin(Join join, QueryInfo queryInfo) {
		List<TupleExpr> joinArgs = new ArrayList<TupleExpr>();
		collectJoinArgs(join, joinArgs);
		return new NJoin(joinArgs, queryInfo);
	}
	
	
	/**
	 * Collect join arguments by descending the query tree (recursively).
	 * 
	 * @param node
	 * @param joinArgs
	 */
	protected static void collectJoinArgs(TupleExpr node, List<TupleExpr> joinArgs) {
		
		if (node instanceof Join) {
			collectJoinArgs(((Join)node).getLeftArg(), joinArgs);
			collectJoinArgs(((Join)node).getRightArg(), joinArgs);
		} else
			joinArgs.add(node);
	}
}
